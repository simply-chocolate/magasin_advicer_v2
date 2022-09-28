package sap_api_wrapper

import (
	"encoding/json"
	"fmt"
)

type SapApiGetStockTransfersResult struct {
	Value []struct {
		DocDate            string `json:"DocDate"`
		DocNum             int    `json:"DocNum"`
		CardCode           string `json:"CardCode"`
		StockTransferLines []struct {
			ItemCode      string      `json:"ItemCode"`
			WarehouseCode string      `json:"WarehouseCode"`
			UoMEntry      json.Number `json:"UoMEntry"`
			UoMCode       string      `json:"UoMCode"`
			Quantity      json.Number `json:"Quantity"`
		} `json:"StockTransferLines"`
	} `json:"value"`
	NextLink string `json:"odata.nextLink"`
}

type SapApiGetStockTransfersReturn struct {
	Body *SapApiGetStockTransfersResult
}

func SapApiGetStockTransfers(params SapApiQueryParams) (SapApiGetStockTransfersReturn, error) {
	client, err := GetSapApiAuthClient()
	if err != nil {
		return SapApiGetStockTransfersReturn{}, err
	}

	resp, err := client.
		R().
		SetResult(SapApiGetStockTransfersResult{}).
		SetQueryParams(params.AsReqParams()).
		Get("StockTransfers")
	if err != nil {
		return SapApiGetStockTransfersReturn{}, err
	}

	if resp.Result() == nil {
		return SapApiGetStockTransfersReturn{}, fmt.Errorf("no advices from this call")
	}

	return SapApiGetStockTransfersReturn{
		Body: resp.Result().(*SapApiGetStockTransfersResult),
	}, nil
}

func SapApiGetStockTransfers_AllPages(params SapApiQueryParams) (SapApiGetStockTransfersReturn, error) {
	res := SapApiGetStockTransfersResult{}
	for page := 0; ; page++ {
		params.Skip = page * 20

		getStockTransferRes, err := SapApiGetStockTransfers(params)
		if err != nil {
			return SapApiGetStockTransfersReturn{}, err
		}

		res.Value = append(res.Value, getStockTransferRes.Body.Value...)

		if getStockTransferRes.Body.NextLink == "" {
			break
		}
	}

	return SapApiGetStockTransfersReturn{
		Body: &res,
	}, nil
}
