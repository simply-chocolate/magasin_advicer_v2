package sap_api_wrapper

import (
	"fmt"
)

type adviceBody struct {
	AdviceStatus string `json:"U_CCF_AdviceStatus"` // N = Default is not sent | Y = Is sent | S = Send again
}

type SapRequestError struct {
	Error struct {
		Code    int `json:"code"`
		Message struct {
			Language     string `json:"lang"`
			ErrorMessage string `json:"value"`
		} `json:"message"`
	} `json:"error"`
}

// Takes the Gs1Status and Gs1 Response and updates the item in SAP
// docEntry : An int telling SAP which document to change
// adviceStats : N = Default is not sent | Y = Is sent | S = Send again
// docType : Order | StockTransfer -> DeliveryNote | StockTransfer
func SetAdviceStatus(docEntry int, adviceStatus string, docType string) error {

	var body adviceBody
	body.AdviceStatus = adviceStatus

	client, err := GetSapApiAuthClient()
	if err != nil {
		fmt.Println("Error getting an authenticaed client")
		return err
	}

	if docType != "DeliveryNotes" && docType != "StockTransfers" {
		return fmt.Errorf("the doctype %v is not known. Valid types are: DeliveryNotes and StockTransfers", docType)
	}
	if docEntry == 0 {
		return fmt.Errorf("the DocEntry cannot be 0..")
	}

	resp, err := client.
		//DevMode().
		R().
		EnableDump().
		SetHeader("Content-Type", "application/json").
		SetErrorResult(SapRequestError{}).
		SetBody(body).
		Patch(fmt.Sprintf("%v(%v)", docType, docEntry))

	if err != nil {
		return err
	}

	if resp.IsErrorState() {
		fmt.Printf("resp is err statusCode: %v. Dump: %v\n", resp.StatusCode, resp.Dump())
		return resp.Err
	}

	return nil
}
