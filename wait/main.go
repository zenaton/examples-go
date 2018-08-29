package main

import (
	_ "github.com/zenaton/examples-go/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflow"
)

func main() {
	workflow.WaitWorkflow.Dispatch()
}
