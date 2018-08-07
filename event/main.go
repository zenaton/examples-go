package main

import (
	"time"

	_ "github.com/zenaton/examples-go/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflow"
)

func main() {

	workflow.EventWorkflow.Dispatch()

	time.Sleep(2 * time.Second)

	workflow.EventWorkflow.WhereID("MyId").Send("MyEvent", nil)
}
