package main

import (
	"fmt"
	"magasin_advicer/teams_notifier"
	"magasin_advicer/utils"

	"time"
)

func main() {
	utils.LoadEnvs()

	fmt.Printf("%v: Started the Script \n", time.Now().UTC().Format("2006-01-02 15:04:05"))

	businessPartners, CardCodesString, err := utils.GetBusinessPartnersFromSap()
	if err != nil {
		teams_notifier.SendUnknownErrorToTeams(err)
		return
	}

	validItems, err := utils.GetValidItemsFromSap()
	if err != nil {
		teams_notifier.SendUnknownErrorToTeams(err)
		return
	}

	//fmt.Printf("%v: Handling Stock Transfers \n", time.Now().UTC().Format("2006-01-02 15:04:05"))
	err = utils.HandleCreateAdviceStockTransfers(businessPartners, CardCodesString, validItems)
	if err != nil {
		teams_notifier.SendUnknownErrorToTeams(err)
	}

	//fmt.Printf("%v: Success handling Stock Transfers \n", time.Now().UTC().Format("2006-01-02 15:04:05"))
	//fmt.Printf("%v: Handling Orders \n", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err = utils.HandleCreateAdviceOrders(businessPartners, CardCodesString, validItems)
	if err != nil {
		teams_notifier.SendUnknownErrorToTeams(err)
	}

	//fmt.Printf("%v: Success handling Orders \n", time.Now().UTC().Format("2006-01-02 15:04:05"))
	fmt.Printf("%v: Success \n", time.Now().UTC().Format("2006-01-02 15:04:05"))

}
