package main

import (
	_ "github.com/zenaton/examples-go/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflow"
)

func main() {

	//todo:
	workflow.WaitEventWorkflow.Dispatch()
}
