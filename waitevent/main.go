package main

import (
	_ "github.com/zenaton/examples-go/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflows"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
	"time"
)

func main() {
	WaitEventWorkflow := workflow.New(&workflows.WaitEvent{})
	WaitEventWorkflow.Dispatch2()

	time.Sleep(2 * time.Second)

	WaitEventWorkflow.WhereID("MyId").Send("MyEvent", nil)
}
