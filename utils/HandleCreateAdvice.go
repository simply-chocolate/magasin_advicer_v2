package utils

import (
	"fmt"
	"magasin_advicer/sap_api_wrapper"
	"strconv"
	"strings"
)

// This function will be called from the main, and call the functions that needs to do stuff, in order to create the advices.
func HandleCreateAdvice() error {
	cardCodes := map[string]string{
		"100084": "10",
		"102024": "15",
		"100087": "20",
		"100334": "25",
		"100089": "30",
		"100088": "40",
		"100090": "50",
		"212868": "60",
	}

	adviceCache, err := ReadAdviceCache()
	if err != nil {
		return err
	}

	stockTransfers, err := sap_api_wrapper.SapApiGetStockTransfers_AllPages(sap_api_wrapper.SapApiQueryParams{
		Select:  []string{"DocDate", "DocNum", "CardCode", "StockTransferLines"},
		OrderBy: []string{"DocNum asc"},
		Filter:  fmt.Sprintf("DocNum gt %v and contains(CardName,'Magasin')", adviceCache.LastAdviceDocNum),
		// Filer: "DocNum eq 81651" // For when we need to create a specific advice..
	})
	if err != nil {
		return err
	}
	if len(stockTransfers.Body.Value) == 0 {
		return fmt.Errorf("no advices from this call")
	}

	sapApiItemsResult, err := sap_api_wrapper.SapApiGetItems_AllPages(sap_api_wrapper.SapApiQueryParams{
		Select: []string{"ItemCode", "ItemBarCodeCollection", "UpdateDate", "UpdateTime"},
		Filter: "Valid eq 'Y'",
	})
	if err != nil {
		fmt.Println("error, something went wrong while getting the items and their properties")
	}

	for _, stockTransfer := range stockTransfers.Body.Value {

		WarehouseCode, cardCodeExists := cardCodes[stockTransfer.CardCode]
		if !cardCodeExists {
			fmt.Printf("cardCode is not a magasin, stockTransfer: %v \n", stockTransfer.DocNum)
			continue
		}
		if WarehouseCode != stockTransfer.StockTransferLines[0].WarehouseCode {
			fmt.Printf("warehouse does not match cardcode, stockTransfer %v \n", stockTransfer.DocNum)
			continue
		}

		res := "\"Følgeseddel\";\"Indkøbsnummer\";\"Stregkode\";\"Indkøbsantal\";\"Hus\""

		for _, stockTransferLine := range stockTransfer.StockTransferLines {
			var barcode string
			for _, items := range sapApiItemsResult.Body.Value {
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
				return fmt.Errorf("error converting quantity to float at stockTransfer: %v", stockTransfer.DocNum)
			}

			if barcode == "" {
				fmt.Printf("This line has no barcode so we just ignore it. Itemnumber:%v and Stockfransfer:%v", stockTransferLine.ItemCode, stockTransfer.DocNum)
				continue
			}

			res += fmt.Sprintf("\n\"%v\";\"Magasin\";\"%s\";\"%v\";\"%s\"", stockTransfer.DocNum, strings.ReplaceAll(barcode, "\"", "\"\""), int(quantityAsInt), strings.ReplaceAll(stockTransferLine.WarehouseCode, "\"", "\"\""))
		}

		SendFileFtp(fmt.Sprintf("%v_Reciept_Magasin_%v.csv", stockTransfer.DocNum, stockTransfer.StockTransferLines[0].WarehouseCode), res)
		adviceCache.LastAdviceDocNum = strconv.Itoa(stockTransfer.DocNum)
	}

	if err = WriteAdviceCache(adviceCache); err != nil {
		return fmt.Errorf("error at stockTransfer: %v adding DocNum to JSON ", adviceCache)
	}

	return nil
}
