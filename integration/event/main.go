package main

import (
	"time"

	"github.com/zenaton/examples-go/integration/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflows"
)

func init() { client.SetEnv("event.env") }
func main() {

	workflows.EventWorkflow.NewInstance().Dispatch()

	time.Sleep(2 * time.Second)

	workflows.EventWorkflow.WhereID("MyId").Send("MyEvent", nil)
}
