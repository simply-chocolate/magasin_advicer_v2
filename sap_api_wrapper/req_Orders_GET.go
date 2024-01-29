package sap_api_wrapper

import (
	"encoding/json"
	"fmt"
)

type SapApiGetOrdersResult struct {
	Value []struct {
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

type SapApiGetOrdersReturn struct {
	Body *SapApiGetOrdersResult
}

func SapApiGetOrders(params SapApiQueryParams) (SapApiGetOrdersReturn, error) {
	client, err := GetSapApiAuthClient()
	if err != nil {
		return SapApiGetOrdersReturn{}, err
	}

	resp, err := client.
		//DevMode().
		R().
		SetErrorResult(SapApiErrorResult{}).
		SetSuccessResult(SapApiGetOrdersResult{}).
		SetQueryParams(params.AsReqParams()).
		Get("DeliveryNotes")
	if err != nil {
		if resp.ErrorResult() == nil {
			return SapApiGetOrdersReturn{}, fmt.Errorf("error getting orders: %v", err)
		}
		return SapApiGetOrdersReturn{}, fmt.Errorf("error getting orders: %v", resp.ErrorResult().(*SapApiErrorResult))
	}

	if resp.SuccessResult() == nil {
		return SapApiGetOrdersReturn{}, nil
	}

	return SapApiGetOrdersReturn{
		Body: resp.SuccessResult().(*SapApiGetOrdersResult),
	}, nil
}

func SapApiGetOrders_AllPages(params SapApiQueryParams) (SapApiGetOrdersReturn, error) {
	res := SapApiGetOrdersResult{}
	for page := 0; ; page++ {
		params.Skip = page * 20

		sapApiGetOrders, err := SapApiGetOrders(params)
		if err != nil {
			return SapApiGetOrdersReturn{}, err
		}

		res.Value = append(res.Value, sapApiGetOrders.Body.Value...)

		if sapApiGetOrders.Body.NextLink == "" {
			break
		}
	}

	return SapApiGetOrdersReturn{
		Body: &res,
	}, nil
}
