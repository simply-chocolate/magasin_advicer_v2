package sap_api_wrapper

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

type SapApiPostLoginResult struct {
	SessionTimeout time.Duration `json:"SessionTimeout"`
}

type SapApiPostLoginReturn struct {
	Cookies []*http.Cookie
	Body    *SapApiPostLoginResult
}

func SapApiPostLogin(optParams ...int) (SapApiPostLoginReturn, error) {
	retries := 0
	if len(optParams) >= 1 {
		retries = optParams[0]
	}

	resp, err := GetSapApiBaseClient().
		R().
		SetResult(SapApiPostLoginResult{}).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"CompanyDB": os.Getenv("SAP_DB_NAME"),
			"UserName":  os.Getenv("SAP_UN"),
			"Password":  os.Getenv("SAP_PW"),
		}).
		Post("Login")
	if err != nil {
		return SapApiPostLoginReturn{}, err
	}

	if resp.IsError() {
		if retries >= 10 {
			return SapApiPostLoginReturn{}, errors.New("failed to login to SAP API")
		}

		fmt.Println("Request failed, retrying...")
		time.Sleep(3 * time.Second)

		return SapApiPostLogin(retries + 1)
	}

	return SapApiPostLoginReturn{
		Cookies: resp.Cookies(),
		Body:    resp.Result().(*SapApiPostLoginResult),
	}, nil
}
