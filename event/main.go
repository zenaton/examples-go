package main

import (
	"time"

	"github.com/twinj/uuid"
	_ "github.com/zenaton/examples-go" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflows"
)

func main() {

	id := uuid.NewV4().String()

	workflows.EventWorkflow.New(id).Dispatch()

	time.Sleep(1 * time.Second)

	workflows.EventWorkflow.WhereID(id).Send("MyEvent", nil)
}
