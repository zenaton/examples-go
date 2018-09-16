package main

import (
	"time"

	"github.com/twinj/uuid"
	_ "github.com/zenaton/examples-go" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflows"
)

func main() {

	id := uuid.NewV4().String()

	workflows.WaitEventWorkflow.New(id).Dispatch()

	time.Sleep(2 * time.Second)

	workflows.WaitEventWorkflow.WhereID(id).Send("MyEvent", "some data")
}
