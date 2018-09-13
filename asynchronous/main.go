package main

import (
	"log"

	// initialize client with credentials
	_ "github.com/zenaton/examples-go/client"
	"github.com/zenaton/examples-go/workflows"
)

func main() {

	err := workflows.AsynchronousWorkflow.NewInstance().Dispatch()

	if err != nil {
		log.Fatal(err)
	}
}
