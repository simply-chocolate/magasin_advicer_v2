package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvs() {
	godotenv.Load()
	if os.Getenv("SAP_DB_NAME") == "" {
		panic("Error loading environment variable SAP_DB_NAME")
	}
	if os.Getenv("SAP_UN") == "" {
		panic("Error loading environment variable SAP_UN")
	}
	if os.Getenv("SAP_PW") == "" {
		panic("Error loading environment variable SAP_PW")
	}
	if os.Getenv("SAP_URL") == "" {
		panic("Error loading environment variable SAP_URL")
	}
	if os.Getenv("TEAMS_WEBHOOK_URL") == "" {
		panic("Error loading environment variable TEAMS_WEBHOOK_URL")
	}
	if os.Getenv("FTP_HOST") == "" {
		panic("Error loading the environment variable FTP_HOST")
	}
	if os.Getenv("FTP_PORT") == "" {
		panic("Error loading the environment variable FTP_PORT")
	}
	if os.Getenv("FTP_USER") == "" {
		panic("Error loading the environment variable FTP_USER")
	}
	if os.Getenv("FTP_PASS") == "" {
		panic("Error loading the environment variable FTP_PASS")
	}
}
