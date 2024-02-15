package sap_api_wrapper

import (
	"fmt"
	"magasin_advicer/teams_notifier"
)

type SapApiGetBusinessPartnersResult struct {
	Value    []BusinessPartner `json:"value"`
	NextLink string            `json:"odata.nextLink"`
}

type BusinessPartner struct {
	CardCode      string `json:"CardCode"`
	CardName      string `json:"CardName"`
	AdviceWhsCode string `json:"U_CCF_Advice_WhsCode"`
}

type SapApiGetBusinessPartnersReturn struct {
	Body *SapApiGetBusinessPartnersResult
}

func SapApiGetBusinessPartners(params SapApiQueryParams) (SapApiGetBusinessPartnersReturn, error) {
	client, err := GetSapApiAuthClient()
	if err != nil {
		fmt.Println("Error getting an authenticaed client: ", err)
		return SapApiGetBusinessPartnersReturn{}, err
	}

	resp, err := client.
		//DevMode().
		R().
		SetSuccessResult(SapApiGetBusinessPartnersResult{}).
		SetErrorResult(SapApiErrorResult{}).
		SetQueryParams(params.AsReqParams()).
		Get("BusinessPartners")
	if err != nil {
		fmt.Println(err)
		return SapApiGetBusinessPartnersReturn{}, err
	}

	if resp.IsErrorState() {
		response := resp.ErrorResult().(*SapApiErrorResult)
		teams_notifier.SendRequestsReturnErrorToTeams("BusinessPartners", "Get", fmt.Sprint(resp), response.Error.Message.Value, "sapAPI")
		return SapApiGetBusinessPartnersReturn{}, fmt.Errorf("error getting orders")
	}

	return SapApiGetBusinessPartnersReturn{
		Body: resp.SuccessResult().(*SapApiGetBusinessPartnersResult),
	}, nil

}

func SapApiGetBusinessPartners_AllPages(params SapApiQueryParams) (SapApiGetBusinessPartnersReturn, error) {
	res := SapApiGetBusinessPartnersResult{}
	for page := 0; ; page++ {
		params.Skip = page * 20

		getItemsRes, err := SapApiGetBusinessPartners(params)
		if err != nil {
			return SapApiGetBusinessPartnersReturn{}, err
		}

		res.Value = append(res.Value, getItemsRes.Body.Value...)

		if getItemsRes.Body.NextLink == "" {
			break
		}
	}

	return SapApiGetBusinessPartnersReturn{
		Body: &res,
	}, nil
}
