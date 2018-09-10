package main

import (
	"time"

	_ "github.com/zenaton/examples-go/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflows"
)

func main() {
	workflows.TestEventValueWorkflow.NewInstance().Dispatch()

	time.Sleep(2 * time.Second)
	workflows.TestEventValueWorkflow.WhereID("TestEventValueID").Send("MyOtherEvent", nil)

	time.Sleep(2 * time.Second)

	workflows.TestEventValueWorkflow.WhereID("TestEventValueID").Send("MyEvent", "test data")
}
