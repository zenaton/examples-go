package client

import (
	"os"

	"github.com/zenaton/zenaton-go/v1/zenaton"
)

func init() {

	// these environment variables are all necessary
	var appID = os.Getenv("ZENATON_APP_ID")
	if appID == "" {
		panic("Please add your Zenaton application id on '.env' file (https://zenaton.com/app/api)")
	}

	var apiToken = os.Getenv("ZENATON_API_TOKEN")
	if apiToken == "" {
		panic("Please add your Zenaton api token on '.env' file (https://zenaton.com/app/api)")
	}

	var appEnv = os.Getenv("ZENATON_APP_ENV")
	if appEnv == "" {
		panic("Please add your Zenaton environment on '.env' file (https://zenaton.com/app/api)")
	}

	// init Zenaton client
	zenaton.InitClient(appID, apiToken, appEnv)
	//todo: get env variables from .env file
}
