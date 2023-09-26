package sap_api_wrapper

type SapApiErrorResult struct {
	Error struct {
		Code    int `json:"code"`
		Message struct {
			Lang  string `json:"lang"`
			Value string `json:"value"`
		} `json:"message"`
	} `json:"error"`
}
