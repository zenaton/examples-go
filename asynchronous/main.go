package main

import (
	// initialize client with credentials
	_ "github.com/zenaton/examples-go"
	"github.com/zenaton/examples-go/workflows"
)

func main() {

	workflows.AsynchronousWorkflow.New().Dispatch()
}
