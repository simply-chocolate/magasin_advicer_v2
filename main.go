package main

import (
	"fmt"
	"log"
	"magasin_advicer/sap_api_wrapper"
	"magasin_advicer/utils"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = utils.HandleCreateAdvice()
	if err != nil {
		fmt.Printf("%v %v \n", time.Now().Format("2006-01-02 15:04:05"), err)
	}

	sap_api_wrapper.SapApiPostLogout()
	fmt.Printf("%v Success \n", time.Now().Format("2006-01-02 15:04:05"))
}
