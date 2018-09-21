package examples_go

import (
	"fmt"

	"os"

	"github.com/joho/godotenv"
	"github.com/zenaton/zenaton-go/v1/zenaton"
	"runtime"
	"strings"
)

func init() {
	SetEnv()
}

func SetEnv() {

	//these steps are necessary as client.go will be called from two different locations, here and in the agent
	// this step will help actually get the path to this file, (not the path to the location in which the application
	_, thisFilePath, _, ok := runtime.Caller(0)
	if !ok {
		panic(thisFilePath)
	}
	thisFilePath = strings.Replace(thisFilePath, "/client.go", "/.env", -1)

	// alternatively, you could just set your env variables in a .bashrc or .bash_profile and not load them from a .env file
	err := godotenv.Load(thisFilePath)
	if err != nil {
		fmt.Println("error: ", err)
		panic("Error loading .env file")
	}

	//make sure that all required environment variables are present
	appID := os.Getenv("ZENATON_APP_ID")
	if appID == "" {
		panic("Please add ZENATON_APP_ID env variable (https://zenaton.com/app/api)")
	}

	apiToken := os.Getenv("ZENATON_API_TOKEN")
	if apiToken == "" {
		panic("Please add ZENATON_API_TOKEN env variable (https://zenaton.com/app/api)")
	}

	appEnv := os.Getenv("ZENATON_APP_ENV")
	if appEnv == "" {
		panic("Please add ZENATON_APP_ENV env variable (https://zenaton.com/app/api)")
	}

	zenaton.InitClient(appID, apiToken, appEnv)
}
