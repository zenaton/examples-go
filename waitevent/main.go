package main

import (
	"time"

	_ "github.com/zenaton/examples-go/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflows"
)

func main() {
	workflows.WaitEventWorkflow.WhereID("MyId").Kill()
	workflows.WaitEventWorkflow.NewInstance().Dispatch()

	time.Sleep(2 * time.Second)

	workflows.WaitEventWorkflow.WhereID("MyId").Send("MyEvent", nil)
}
