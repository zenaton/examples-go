package main

import (
	"time"

	_ "github.com/zenaton/examples-go/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflows"
)

func main() {

	workflows.EventWorkflow.NewInstance().Dispatch()

	time.Sleep(2 * time.Second)

	workflows.SequentialWorkflow.WhereID("MyId").Find()

	workflows.EventWorkflow.WhereID("MyId").Send("MyEvent", nil)
}
