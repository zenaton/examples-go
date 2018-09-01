package main

import (
	"time"

	_ "github.com/zenaton/examples-go/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflows"
)

func main() {
	WaitEventWorkflow := workflows.WaitEventWorkflow()
	WaitEventWorkflow.Dispatch()

	time.Sleep(2 * time.Second)

	WaitEventWorkflow.WhereID("MyId").Send("MyEvent", nil)
}
