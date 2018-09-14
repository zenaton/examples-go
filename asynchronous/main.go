package main

import (
	// initialize client with credentials
	_ "github.com/zenaton/examples-go/client"
	"github.com/zenaton/examples-go/workflows"
)

func main() {

	workflows.AsynchronousWorkflow.NewInstance().Dispatch()
}
