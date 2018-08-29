package client

import (
	"os"

	"github.com/zenaton/zenaton-go/v1/zenaton"
)

func init() {
	SetEnv()
}

func SetEnv() {
	var appID = os.Getenv("ZENATON_APP_ID")
	if appID == "" {
		panic("Please add ZENATON_APP_ID env variable (https://zenaton.com/app/api)")
	}

	var apiToken = os.Getenv("ZENATON_API_TOKEN")
	if apiToken == "" {
		panic("Please add ZENATON_API_TOKEN env variable (https://zenaton.com/app/api)")
	}

	var appEnv = os.Getenv("ZENATON_APP_ENV")
	if appEnv == "" {
		panic("Please add ZENATON_APP_ENV env variable(https://zenaton.com/app/api)")
	}

	zenaton.InitClient(appID, apiToken, appEnv)
}
