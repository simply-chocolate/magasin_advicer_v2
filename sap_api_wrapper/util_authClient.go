package sap_api_wrapper

import (
	"net/http"
	"time"

	"github.com/imroc/req/v3"
)

var authCookiesCache []*http.Cookie
var cacheExpiresAt time.Time

func GetSapApiAuthClient() (*req.Client, error) {
	if authCookiesCache == nil || cacheExpiresAt.Before(time.Now()) {
		loginRes, err := SapApiPostLogin()
		if err != nil {
			return nil, err
		}

		authCookiesCache = loginRes.Cookies
		cacheExpiresAt = time.Now().Add((loginRes.Body.SessionTimeout - 1) * time.Minute)
	}

	return GetSapApiBaseClient().SetCommonCookies(authCookiesCache...), nil
}
