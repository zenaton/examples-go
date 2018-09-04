package client

import (
	"fmt"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/zenaton/zenaton-go/v1/zenaton"
	"strings"
)

func init() {
	SetEnv()
}

func SetEnv() {

	_, thisFilePath, _, ok := runtime.Caller(0)
	if !ok {
		panic(thisFilePath)
	}

	thisFilePath = strings.Replace(thisFilePath, "/client.go", "/.env", -1)
	variables, err := godotenv.Read(thisFilePath)
	if err != nil {
		fmt.Println("error: ", err)
		panic("Error loading .env file")
	}

	//make sure that all required environment variables are present
	appID, ok := variables["ZENATON_APP_ID"]
	if !ok {
		panic("Please add ZENATON_APP_ID env variable (https://zenaton.com/app/api)")
	}

	apiToken, ok := variables["ZENATON_API_TOKEN"]
	if !ok {
		panic("Please add ZENATON_API_TOKEN env variable (https://zenaton.com/app/api)")
	}

	appEnv, ok := variables["ZENATON_APP_ENV"]
	if !ok {
		panic("Please add ZENATON_APP_ENV env variable(https://zenaton.com/app/api)")
	}

	zenaton.InitClient(appID, apiToken, appEnv)
}
