package main

import (
	"fmt"
	"magasin_advicer/teams_notifier"
	"magasin_advicer/utils"

	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	utils.LoadEnvs()

	fmt.Printf("%v: Started the Script \n", time.Now().UTC().Format("2006-01-02 15:04:05"))

	fmt.Printf("%v: Handling Stock Transfers \n", time.Now().UTC().Format("2006-01-02 15:04:05"))
	err := utils.HandleCreateAdviceStockTransfers()
	if err != nil {
		teams_notifier.SendUnknownErrorToTeams(err)
	}
	fmt.Printf("%v: Success handling Stock Transfers \n", time.Now().UTC().Format("2006-01-02 15:04:05"))

	fmt.Printf("%v: Handling Orders \n", time.Now().UTC().Format("2006-01-02 15:04:05"))
	err = utils.HandleCreateAdviceOrders()
	if err != nil {
		teams_notifier.SendUnknownErrorToTeams(err)
	}
	fmt.Printf("%v: Success handling Orders \n", time.Now().UTC().Format("2006-01-02 15:04:05"))

	fmt.Printf("%v: Started the Cron Scheduler", time.Now().UTC().Format("2006-01-02 15:04:05"))

	s := gocron.NewScheduler(time.UTC)
	_, _ = s.Cron("0 4,15 * * 1-5").SingletonMode().Do(func() {

		fmt.Printf("%v: Handling Stock Transfers \n", time.Now().UTC().Format("2006-01-02 15:04:05"))
		err = utils.HandleCreateAdviceStockTransfers()
		if err != nil {
			teams_notifier.SendUnknownErrorToTeams(err)
		}
		fmt.Printf("%v: Success handling Stock Transfers \n", time.Now().UTC().Format("2006-01-02 15:04:05"))

		fmt.Printf("%v: Handling Orders \n", time.Now().UTC().Format("2006-01-02 15:04:05"))
		err = utils.HandleCreateAdviceOrders()
		if err != nil {
			teams_notifier.SendUnknownErrorToTeams(err)
		}
		fmt.Printf("%v: Success handling Orders \n", time.Now().UTC().Format("2006-01-02 15:04:05"))

	})

	s.StartBlocking()

}
