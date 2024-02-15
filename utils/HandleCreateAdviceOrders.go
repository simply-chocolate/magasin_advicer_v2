package utils

import (
	"fmt"
	"magasin_advicer/sap_api_wrapper"
	"magasin_advicer/teams_notifier"
	"os"
	"strings"
)

// TODO: Tilføj i filtre at dem med "N" IKKE skal komme med?
// TODO: Tjek leveringer : Advice: 111825, 111826
// TODO: Det ligner at det gik galt specifikt i ordre modulet, at den blev ved dér, så undersøg om der kan være noget her som gør at den ikke bliver markeret korrekt i SAP
// TODO: Undersøg hvilken ordren den TOMME advis kommer fra og hvordan det kan være?

// This function will be called from the main, and call the functions that needs to do stuff, in order to create the advices.
// DeliveryNotes are created from Orders - Naming is bit of both right now.
func HandleCreateAdviceOrders(
	businessPartners map[string]sap_api_wrapper.BusinessPartner,
	cardCodeString string,
	validItemsSap sap_api_wrapper.SapApiGetItemsReturn) error {

	orders, err := GetSapDeliveryNotes(cardCodeString)
	if err != nil {
		return err
	}
	if len(orders.Body.Value) == 0 {
		return nil
	}

	var magasinAdvicesInfo []teams_notifier.MagasinAdviceInfo
	for _, order := range orders.Body.Value {
		docNum := fmt.Sprint(order.DocNum)
		businessPartner, exists := businessPartners[order.CardCode]
		if !exists {
			fmt.Printf("CardCode: %v does not exists in our cardcodes?", order.CardCode)
			continue
		}
		warehouseCode := businessPartner.AdviceWhsCode

		if len(order.OrderLines) == 0 {
			teams_notifier.SendUnknownErrorToTeams(fmt.Errorf("DeliveryNote %v does not have any lines but made it through the other criteria", order.DocNum))
			continue
		}

		orderNumber := order.OrderNumber
		if order.OrderNumber == "" {
			orderNumber = "Magasin"
		}

		res := "\"Følgeseddel\";\"Indkøbsnummer\";\"Stregkode\";\"Indkøbsantal\";\"Hus\""

		for _, orderLine := range order.OrderLines {
			if orderLine.UnitsOfMeasurement == "" {
				return fmt.Errorf("error at order: %v. Error: UnitsOfMeasurement is undefined", order.DocNum)
			}

			var barcode string
			// TODO: This could do with a rewamp. But it works.
			for _, items := range validItemsSap.Body.Value {
				if items.ItemCode == orderLine.ItemCode {
					for _, barCodeColletion := range items.ItemBarCodeCollection {
						if barCodeColletion.UoMEntry == "1" {
							barcode = barCodeColletion.Barcode
						}
					}
				}
			}

			unitsPerQuantityAsFloat, err := orderLine.UnitsOfMeasurement.Float64()
			if err != nil {
				return fmt.Errorf("error converting unitsPerQuantity to float at order: %v. Error:%v", order.DocNum, err)
			}

			quantityAsFloat, err := orderLine.Quantity.Float64()
			if err != nil {
				return fmt.Errorf("error converting quantity to float at order: %v. Error:%v", order.DocNum, err)
			}

			quantity := quantityAsFloat * unitsPerQuantityAsFloat

			if barcode == "" {
				continue // This line has no barcode so we just ignore it.
			}

			res += fmt.Sprintf("\n\"%v\";\"%v\";\"%s\";\"%v\";\"%s\"", docNum, strings.ReplaceAll(orderNumber, "\"", "\"\""), strings.ReplaceAll(barcode, "\"", "\"\""), int(quantity), strings.ReplaceAll(warehouseCode, "\"", "\"\""))
		}

		if os.Getenv("DEVMODE") == "false" {
			err = SendFileFtp(fmt.Sprintf("%v_Order_Reciept_Magasin_%v.csv", docNum, warehouseCode), res, "MAGASIN")
			if err != nil {
				return fmt.Errorf("error sending DeliveryNote %v to FTP: %v", docNum, err)
			}
		} else {
			err = SaveDataAsCSV(fmt.Sprintf("%v_Order_Reciept_Magasin_%v.csv", docNum, warehouseCode), res, "MAGASIN")
			if err != nil {
				return fmt.Errorf("error saving DeliveryNote %v to CSV: %v", docNum, err)
			}
		}

		err = sap_api_wrapper.SetAdviceStatus(order.DocEntry, "Y", "DeliveryNotes")
		if err != nil {
			teams_notifier.SendUnknownErrorToTeams(fmt.Errorf("error changing advice status to 'Y' for delivery note: %v. \n error:%v", docNum, err))
			// TODO: Måske skal vi bruge vores Cache til dette istedet or smide DocNum ind på dem der fejler?
		}

		var magasinAdviceInfo teams_notifier.MagasinAdviceInfo
		magasinAdviceInfo.AdviceNumber = order.DocNum
		magasinAdviceInfo.HouseNumber = warehouseCode
		magasinAdvicesInfo = append(magasinAdvicesInfo, magasinAdviceInfo)
	}

	teams_notifier.SendAdviceSuccesToTeams(magasinAdvicesInfo, "MAGASIN: Order")
	return nil
}
