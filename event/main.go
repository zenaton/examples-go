package main

import (
	"time"

	_ "github.com/zenaton/examples-go/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflows"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

func main() {

	EventWorkflow := workflow.New(&workflows.Event{})
	EventWorkflow.Dispatch2()

	time.Sleep(2 * time.Second)

	EventWorkflow.WhereID("MyId").Send("MyEvent", nil)
}
