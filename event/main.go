package main

import (
	"time"

	_ "github.com/zenaton/examples-go/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflows"
)

func main() {

	EventWorkflow := workflows.EventWorkflow()
	EventWorkflow.Dispatch()

	time.Sleep(2 * time.Second)

	EventWorkflow.WhereID("MyId").Send("MyEvent", nil)
}
