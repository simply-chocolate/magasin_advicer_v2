package utils

import (
	"fmt"
	"magasin_advicer/sap_api_wrapper"
	"magasin_advicer/teams_notifier"
	"os"
	"strings"
)

// This function will be called from the main, and call the functions that needs to do stuff, in order to create the advices.
func HandleCreateAdviceStockTransfers(
	businessPartners map[string]sap_api_wrapper.BusinessPartner,
	cardCodeString string,
	validItemsSap sap_api_wrapper.SapApiGetItemsReturn) error {

	stockTransfers, err := GetSapStockTransfers(cardCodeString)
	if err != nil {
		return err
	}
	if len(stockTransfers.Body.Value) == 0 {
		return nil
	}

	var magasinAdvicesInfo []teams_notifier.MagasinAdviceInfo
	for _, stockTransfer := range stockTransfers.Body.Value {
		docNum := fmt.Sprint(stockTransfer.DocNum)

		businessPartner, exists := businessPartners[stockTransfer.CardCode]
		if !exists {
			fmt.Printf("CardCode: %v does not exists in our cardcodes?", stockTransfer.CardCode)
			continue
		}
		WarehouseCode := businessPartner.AdviceWhsCode

		if len(stockTransfer.StockTransferLines) == 0 {
			teams_notifier.SendUnknownErrorToTeams(fmt.Errorf("StockTransfer %v does not have any lines but made it through the other criteria", stockTransfer.DocNum))
			continue
		}
		if WarehouseCode != stockTransfer.StockTransferLines[0].WarehouseCode {
			continue // Warehouse does not match expected warehouse
		}

		res := "\"Følgeseddel\";\"Indkøbsnummer\";\"Stregkode\";\"Indkøbsantal\";\"Hus\""

		// TODO: This could do with a rewamp. But it works.
		for _, stockTransferLine := range stockTransfer.StockTransferLines {
			var barcode string
			for _, items := range validItemsSap.Body.Value {
				if items.ItemCode == stockTransferLine.ItemCode {
					for _, barCodeColletion := range items.ItemBarCodeCollection {
						if barCodeColletion.UoMEntry == stockTransferLine.UoMEntry {
							barcode = barCodeColletion.Barcode
						}
					}
				}
			}

			quantityAsFloat, err := stockTransferLine.Quantity.Float64()
			if err != nil {
				return fmt.Errorf("error converting quantity to float at stockTransfer: %v. Error:%v", stockTransfer.DocNum, err)
			}

			if barcode == "" {
				continue // This line has no barcode so we just ignore it.
			}

			res += fmt.Sprintf("\n\"%v\";\"Magasin\";\"%s\";\"%v\";\"%s\"", docNum, strings.ReplaceAll(barcode, "\"", "\"\""), int(quantityAsFloat), strings.ReplaceAll(stockTransferLine.WarehouseCode, "\"", "\"\""))
		}

		if os.Getenv("DEVMODE") == "false" {
			err = SendFileFtp(fmt.Sprintf("%v_StockTransfer_Reciept_Simply_%v.csv", docNum, stockTransfer.StockTransferLines[0].WarehouseCode), res, "SIMPLY")
			if err != nil {
				return fmt.Errorf("error sending StockTransfer %v to FTP: %v", docNum, err)
			}
		} else {
			err = SaveDataAsCSV(fmt.Sprintf("%v_StockTransfer_Reciept_Magasin_%v.csv", docNum, stockTransfer.StockTransferLines[0].WarehouseCode), res, "MAGASIN")
			if err != nil {
				return fmt.Errorf("error saving StockTransfer %v to CSV: %v", docNum, err)
			}
		}

		err = sap_api_wrapper.SetAdviceStatus(stockTransfer.DocEntry, "Y", "StockTransfers")
		if err != nil {
			teams_notifier.SendUnknownErrorToTeams(fmt.Errorf("error changing advice status to 'Y' for stockTransfer: %v. \n error:%v", docNum, err))
			// TODO: Måske skal vi bruge vores Cache til dette istedet or smide DocNum ind på dem der fejler?
		}

		var magasinAdviceInfo teams_notifier.MagasinAdviceInfo

		magasinAdviceInfo.AdviceNumber = stockTransfer.DocNum
		magasinAdviceInfo.HouseNumber = stockTransfer.StockTransferLines[0].WarehouseCode
		magasinAdvicesInfo = append(magasinAdvicesInfo, magasinAdviceInfo)
	}

	teams_notifier.SendAdviceSuccesToTeams(magasinAdvicesInfo, "SIMPLY: StockTransfers")
	return nil
}
