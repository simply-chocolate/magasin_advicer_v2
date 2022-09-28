package sap_api_wrapper

import (
	"os"

	"github.com/imroc/req/v3"
)

func GetSapApiBaseClient() *req.Client {
	return req.C().SetBaseURL(os.Getenv("SAP_URL"))
}
