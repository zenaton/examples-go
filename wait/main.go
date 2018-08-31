package main

import (
	_ "github.com/zenaton/examples-go/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflows"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

func main() {
	workflow.New(&workflows.Wait{}).Dispatch2()
}
