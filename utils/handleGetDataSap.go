package utils

import (
	"fmt"
	"magasin_advicer/sap_api_wrapper"
)

func GetBusinessPartnersFromSap() (map[string]sap_api_wrapper.BusinessPartner, string, error) {
	resp, err := sap_api_wrapper.SapApiGetBusinessPartners_AllPages(sap_api_wrapper.SapApiQueryParams{
		Select: []string{
			"CardCode",
			"CardName",
			"U_CCF_Advice_WhsCode",
		},
		Filter: "U_CCF_Advice_WhsCode ne null",
	})
	if err != nil {
		return map[string]sap_api_wrapper.BusinessPartner{}, "", err
	}

	BusinessPartners := make(map[string]sap_api_wrapper.BusinessPartner)
	CardCodesString := "("
	for _, businessPartner := range resp.Body.Value {
		BusinessPartners[businessPartner.CardCode] = businessPartner
		CardCodesString += fmt.Sprintf("CardCode eq '%s' or ", businessPartner.CardCode)
	}
	CardCodesString = CardCodesString[:len(CardCodesString)-4]
	CardCodesString += ")"

	return BusinessPartners, CardCodesString, nil
}

func GetValidItemsFromSap() (sap_api_wrapper.SapApiGetItemsReturn, error) {
	validItemsSap, err := sap_api_wrapper.SapApiGetItems_AllPages(sap_api_wrapper.SapApiQueryParams{
		Select: []string{"ItemCode", "ItemBarCodeCollection", "UpdateDate", "UpdateTime"},
		Filter: "Valid eq 'Y'",
	})

	if err != nil {
		return sap_api_wrapper.SapApiGetItemsReturn{}, fmt.Errorf("error getting valid items from SAP")
	}

	return validItemsSap, nil
}

func GetSapDeliveryNotes(CardCodesString string) (sap_api_wrapper.SapApiGetDerliveryNotesReturn, error) {
	lastDateForCache := "2024-02-01"
	deliveryNotes, err := sap_api_wrapper.SapApiGetDeliveryNotes_AllPages(sap_api_wrapper.SapApiQueryParams{
		Select:  []string{"DocEntry", "DocDate", "DocNum", "CardCode", "NumAtCard", "U_CCF_AdviceStatus", "DocumentLines"},
		OrderBy: []string{"DocNum asc"},
		Filter:  fmt.Sprintf("DocDate ge %v and U_CCF_AdviceStatus ne 'Y' and %v", lastDateForCache, CardCodesString),
		//Filter: "DocNum eq 102987", // For when we need to create a specific advice...............
	})

	if err != nil {
		return sap_api_wrapper.SapApiGetDerliveryNotesReturn{}, fmt.Errorf("error getting delivery notes from SAP: %v", err)
	}

	return deliveryNotes, nil
}
func GetSapStockTransfers(CardCodesString string) (sap_api_wrapper.SapApiGetStockTransfersReturn, error) {
	lastDateForCache := "2024-02-01"
	stockTransfers, err := sap_api_wrapper.SapApiGetStockTransfers_AllPages(sap_api_wrapper.SapApiQueryParams{
		Select:  []string{"DocEntry", "DocDate", "DocNum", "CardCode", "U_CCF_AdviceStatus", "StockTransferLines"},
		OrderBy: []string{"DocNum asc"},
		Filter:  fmt.Sprintf("DocDate ge %v and U_CCF_AdviceStatus ne 'Y' and %v", lastDateForCache, CardCodesString),
		//Filter: "DocNum eq 102987", // For when we need to create a specific advice...............
	})

	if err != nil {
		return sap_api_wrapper.SapApiGetStockTransfersReturn{}, fmt.Errorf("error getting stock transfers from SAP: %v", err)
	}

	return stockTransfers, nil
}
