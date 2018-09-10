package main

import (
	"time"

	"github.com/zenaton/examples-go/integration/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflows"
)

func init() { client.SetEnv("waitevent.env") }
func main() {

	workflows.WaitEventWorkflow.NewInstance().Dispatch()

	time.Sleep(2 * time.Second)

	workflows.WaitEventWorkflow.WhereID("MyId").Send("MyEvent", "some data")
}
