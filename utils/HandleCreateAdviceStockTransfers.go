package utils

import (
	"fmt"
	"magasin_advicer/sap_api_wrapper"
	"magasin_advicer/teams_notifier"
	"strconv"
	"strings"
)

// This function will be called from the main, and call the functions that needs to do stuff, in order to create the advices.
func HandleCreateAdviceStockTransfers() error {
	stockTransferCardCodes := map[string]string{
		"100084": "10",
		"100085": "10",
		"102024": "15",
		"100087": "20",
		"100334": "25",
		"100089": "30",
		"100088": "40",
		"100090": "50",
		"212868": "60",
	}

	adviceCache, err := ReadAdviceCache("stockTransfers")
	if err != nil {
		return err
	}

	stockTransfers, err := sap_api_wrapper.SapApiGetStockTransfers_AllPages(sap_api_wrapper.SapApiQueryParams{
		Select:  []string{"DocEntry", "DocDate", "DocNum", "CardCode", "U_CCF_AdviceStatus", "StockTransferLines"},
		OrderBy: []string{"DocNum asc"},
		Filter:  fmt.Sprintf("(DocNum gt %v or U_CCF_AdviceStatus eq 'S') and contains(CardName,'Magasin ')", adviceCache.LastAdviceDocNum),
		//Filter: "DocNum eq 102987", // For when we need to create a specific advice...............
	})

	if err != nil {
		teams_notifier.SendRequestsReturnErrorToTeams("SapApiGetStockTransfers_AllPages", "GET", "Error", err.Error(), "SAP API")
		return nil
	}
	if len(stockTransfers.Body.Value) == 0 {
		teams_notifier.SendNoAdviceToTeams("SIMPLY: StockTransfers")
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
	for _, stockTransfer := range stockTransfers.Body.Value {
		var docNum string
		WarehouseCode, cardCodeExists := stockTransferCardCodes[stockTransfer.CardCode]
		if !cardCodeExists {
			continue // CardCode is not a magasin
		}
		if WarehouseCode != stockTransfer.StockTransferLines[0].WarehouseCode {
			continue // Warehouse does not match expected warehouse
		}

		res := "\"Følgeseddel\";\"Indkøbsnummer\";\"Stregkode\";\"Indkøbsantal\";\"Hus\""

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

			quantityAsInt, err := stockTransferLine.Quantity.Float64()
			if err != nil {
				return fmt.Errorf("error converting quantity to float at stockTransfer: %v. Error:%v", stockTransfer.DocNum, err)
			}

			if barcode == "" {
				continue // This line has no barcode so we just ignore it.
			}

			docNum = fmt.Sprint(stockTransfer.DocNum)

			if stockTransfer.AdviceStatus == "S" {
			} else {
				adviceCache.LastAdviceDocNum = strconv.Itoa(stockTransfer.DocNum)
			}

			res += fmt.Sprintf("\n\"%v\";\"Magasin\";\"%s\";\"%v\";\"%s\"", docNum, strings.ReplaceAll(barcode, "\"", "\"\""), int(quantityAsInt), strings.ReplaceAll(stockTransferLine.WarehouseCode, "\"", "\"\""))
		}

		SendFileFtp(fmt.Sprintf("%v_StockTransfer_Reciept_Simply_%v.csv", docNum, stockTransfer.StockTransferLines[0].WarehouseCode), res, "SIMPLY")

		err = sap_api_wrapper.SetAdviceStatus(stockTransfer.DocEntry, "Y", "StockTransfers")
		if err != nil {
			fmt.Println(err)
		}
		var magasinAdviceInfo teams_notifier.MagasinAdviceInfo

		magasinAdviceInfo.AdviceNumber = stockTransfer.DocNum
		magasinAdviceInfo.HouseNumber = stockTransfer.StockTransferLines[0].WarehouseCode
		magasinAdvicesInfo = append(magasinAdvicesInfo, magasinAdviceInfo)
	}

	if adviceCache.LastAdviceDocNum != "" {
		if err = WriteAdviceCache(adviceCache, "orders"); err != nil {
			return fmt.Errorf("error at order: %v adding DocNum to JSON ", adviceCache)
		}
	}

	if len(magasinAdvicesInfo) == 0 {
		fmt.Printf("AdvicesInfo is empty. StockTransfer: %v", stockTransfers.Body.Value)
	}
	teams_notifier.SendAdviceSuccesToTeams(magasinAdvicesInfo, "SIMPLY: StockTransfers")
	return nil
}
