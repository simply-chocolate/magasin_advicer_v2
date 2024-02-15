package sap_api_wrapper

import (
	"encoding/json"
	"fmt"
)

type SapApiGetDeliveryNotesResult struct {
	Value []struct {
		DocEntry     int    `json:"DocEntry"`
		DocDate      string `json:"DocDate"`
		DocNum       int    `json:"DocNum"`
		CardCode     string `json:"CardCode"`
		OrderNumber  string `json:"NumAtCard"`
		AdviceStatus string `json:"U_CCF_AdviceStatus"` // N = Default is not sent | Y = Is sent | S = Send again
		OrderLines   []struct {
			ItemCode           string      `json:"ItemCode"`
			WarehouseCode      string      `json:"WarehouseCode"`
			UoMEntry           json.Number `json:"UoMEntry"`
			UoMCode            string      `json:"UoMCode"`
			Quantity           json.Number `json:"Quantity"`
			UnitsOfMeasurement json.Number `json:"UnitsOfMeasurment"`
		} `json:"DocumentLines"`
	} `json:"value"`
	NextLink string `json:"odata.nextLink"`
}

type SapApiGetDerliveryNotesReturn struct {
	Body *SapApiGetDeliveryNotesResult
}

func SapApiGetDeliveryNotes(params SapApiQueryParams) (SapApiGetDerliveryNotesReturn, error) {
	client, err := GetSapApiAuthClient()
	if err != nil {
		return SapApiGetDerliveryNotesReturn{}, err
	}

	resp, err := client.
		//DevMode().
		R().
		SetErrorResult(SapApiErrorResult{}).
		SetSuccessResult(SapApiGetDeliveryNotesResult{}).
		SetQueryParams(params.AsReqParams()).
		Get("DeliveryNotes")

	if err != nil {
		if resp.ErrorResult() == nil {
			return SapApiGetDerliveryNotesReturn{}, fmt.Errorf("error getting deliveryNotes: %v", err)
		}
		return SapApiGetDerliveryNotesReturn{}, fmt.Errorf("error getting deliveryNotes: %v", resp.ErrorResult().(*SapApiErrorResult))
	}

	if resp.SuccessResult() == nil {
		fmt.Println(resp)
		return SapApiGetDerliveryNotesReturn{}, fmt.Errorf("no deliveryNotes were found")
	}

	return SapApiGetDerliveryNotesReturn{
		Body: resp.SuccessResult().(*SapApiGetDeliveryNotesResult),
	}, nil
}

func SapApiGetDeliveryNotes_AllPages(params SapApiQueryParams) (SapApiGetDerliveryNotesReturn, error) {
	res := SapApiGetDeliveryNotesResult{}
	for page := 0; ; page++ {
		params.Skip = page * 20

		sapApiGetDeliveryNotes, err := SapApiGetDeliveryNotes(params)
		if err != nil {
			return SapApiGetDerliveryNotesReturn{}, err
		}

		res.Value = append(res.Value, sapApiGetDeliveryNotes.Body.Value...)

		if sapApiGetDeliveryNotes.Body.NextLink == "" {
			break
		}
	}

	return SapApiGetDerliveryNotesReturn{
		Body: &res,
	}, nil
}
