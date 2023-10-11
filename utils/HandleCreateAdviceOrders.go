package utils

import (
	"fmt"
	"magasin_advicer/sap_api_wrapper"
	"magasin_advicer/teams_notifier"
	"strconv"
	"strings"
)

// This function will be called from the main, and call the functions that needs to do stuff, in order to create the advices.
func HandleCreateAdviceOrders() error {
	orderCardCodes := map[string]string{
		"100085": "10",
		"102024": "15",
		"100087": "20",
		"100334": "25",
		"100089": "30",
		"100088": "40",
		"100090": "50",
		"212868": "60",
	}

	adviceCache, err := ReadAdviceCache("orders")
	if err != nil {
		return err
	}

	orders, err := sap_api_wrapper.SapApiGetOrders_AllPages(sap_api_wrapper.SapApiQueryParams{
		Select:  []string{"DocDate", "DocNum", "CardCode", "NumAtCard", "DocumentLines"},
		OrderBy: []string{"DocNum asc"},
		Filter:  fmt.Sprintf("DocNum gt %v and startswith(CardName,'Magasin ') and CardCode ne '100084'", adviceCache.LastAdviceDocNum),
		//Filter: "DocNum eq 102987", // For when we need to create a specific advice...............
	})

	if err != nil {
		teams_notifier.SendRequestsReturnErrorToTeams("SapApiGetOrders_AllPages", "GET", "Error", err.Error(), "SAP API")
		return nil
	}
	if len(orders.Body.Value) == 0 {
		teams_notifier.SendNoAdviceToTeams("MAGASIN: Order")
		return nil
	}

	validItemsSap, err := sap_api_wrapper.SapApiGetItems_AllPages(sap_api_wrapper.SapApiQueryParams{
		Select: []string{"ItemCode", "ItemBarCodeCollection", "UpdateDate", "UpdateTime"},
		Filter: "Valid eq 'Y'",
	})
	if err != nil {
		teams_notifier.SendRequestsReturnErrorToTeams("SapApiGetItems_AllPages", "GET", "Error", err.Error(), "SAP API")
		return nil
	}

	var magasinAdvicesInfo []teams_notifier.MagasinAdviceInfo
	for _, order := range orders.Body.Value {

		warehouseCode, cardCodeExists := orderCardCodes[order.CardCode]
		if !cardCodeExists {
			continue // CardCode is not the correct Magasin for orders
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

			res += fmt.Sprintf("\n\"%v\";\"%v\";\"%s\";\"%v\";\"%s\"", order.DocNum, strings.ReplaceAll(orderNumber, "\"", "\"\""), strings.ReplaceAll(barcode, "\"", "\"\""), int(quantity), strings.ReplaceAll(warehouseCode, "\"", "\"\""))
		}

		err = SendFileFtp(fmt.Sprintf("%v_Reciept_Magasin_%v.csv", order.DocNum, warehouseCode), res, "MAGASIN")
		if err != nil {
			teams_notifier.SendRequestsReturnErrorToTeams("SendFileFtp", "POST", "Error", err.Error(), "FTP")
			return nil
		}
		adviceCache.LastAdviceDocNum = strconv.Itoa(order.DocNum)

		var magasinAdviceInfo teams_notifier.MagasinAdviceInfo
		magasinAdviceInfo.AdviceNumber = order.DocNum
		magasinAdviceInfo.HouseNumber = warehouseCode
		magasinAdvicesInfo = append(magasinAdvicesInfo, magasinAdviceInfo)
	}

	if err = WriteAdviceCache(adviceCache, "orders"); err != nil {
		return fmt.Errorf("error at order: %v adding DocNum to JSON ", adviceCache)
	}

	teams_notifier.SendAdviceSuccesToTeams(magasinAdvicesInfo, "MAGASIN: Order")
	return nil
}
