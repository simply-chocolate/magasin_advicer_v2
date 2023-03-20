package main

import (
	"fmt"
	"magasin_advicer/sap_api_wrapper"
	"magasin_advicer/teams_notifier"
	"magasin_advicer/utils"

	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	utils.LoadEnvs()
	fmt.Println("Started the Cron Scheduler for magasin advicer")
	fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05"))

	s := gocron.NewScheduler(time.UTC)
	_, _ = s.Cron("0 4 * * *").SingletonMode().Do(func() {
		fmt.Printf("%v Started the Script \n", time.Now().UTC().Format("2006-01-02 15:04:05"))

		err := utils.HandleCreateAdvice()
		if err != nil {
			teams_notifier.SendUnknownErrorToTeams(err)
		}

		sap_api_wrapper.SapApiPostLogout()
		fmt.Printf("%v Success \n", time.Now().UTC().Format("2006-01-02 15:04:05"))
	})

	s.StartBlocking()

}
