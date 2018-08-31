package main

import (
	_ "github.com/zenaton/examples-go/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflows"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

func main() {
	//workflow.SequentialWorkflow2("bob").Dispatch2()
	workflow.New(&workflows.Sequential{"bob"}).Dispatch2()
}
