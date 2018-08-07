package main

import (
	_ "github.com/zenaton/examples-go/client"
	"github.com/zenaton/examples-go/workflow"
) // initialize client with credentials

func main() {
	workflow.AsynchronousWorkflow.Dispatch()
}
