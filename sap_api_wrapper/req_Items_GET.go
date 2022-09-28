package sap_api_wrapper

import (
	"encoding/json"
)

type SapApiGetItemsResult struct {
	Value []struct {
		ItemCode              string `json:"ItemCode"`
		UpdateDate            string `json:"UpdateDate"`
		UpdateTime            string `json:"UpdateTime"`
		ItemBarCodeCollection []struct {
			UoMEntry json.Number `json:"UoMEntry"`
			Barcode  string      `json:"Barcode"`
		} `json:"ItemBarCodeCollection"`
	} `json:"value"`
	NextLink string `json:"odata.nextLink"`
}

type SapApiGetItemsReturn struct {
	Body *SapApiGetItemsResult
}

func SapApiGetItems(params SapApiQueryParams) (SapApiGetItemsReturn, error) {
	client, err := GetSapApiAuthClient()
	if err != nil {
		return SapApiGetItemsReturn{}, err
	}

	resp, err := client.
		R().
		SetResult(SapApiGetItemsResult{}).
		SetQueryParams(params.AsReqParams()).
		Get("Items")
	if err != nil {
		return SapApiGetItemsReturn{}, err
	}

	return SapApiGetItemsReturn{
		Body: resp.Result().(*SapApiGetItemsResult),
	}, nil
}

func SapApiGetItems_AllPages(params SapApiQueryParams) (SapApiGetItemsReturn, error) {
	res := SapApiGetItemsResult{}
	for page := 0; ; page++ {
		params.Skip = page * 20

		getItemsRes, err := SapApiGetItems(params)
		if err != nil {
			return SapApiGetItemsReturn{}, err
		}

		res.Value = append(res.Value, getItemsRes.Body.Value...)

		if getItemsRes.Body.NextLink == "" {
			break
		}
	}

	return SapApiGetItemsReturn{
		Body: &res,
	}, nil
}
