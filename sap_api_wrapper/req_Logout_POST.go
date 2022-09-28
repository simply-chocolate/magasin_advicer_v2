package sap_api_wrapper

func SapApiPostLogout() error {
	client, err := GetSapApiAuthClient()
	if err != nil {
		println("Error getting auth client for 'Logout'")
		return err
	}

	_, err = client.
		R().
		Post("Logout")
	if err != nil {
		return err
	}

	return nil
}
