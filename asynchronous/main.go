package main

import (
	_ "github.com/zenaton/examples-go/client"
	"github.com/zenaton/examples-go/workflows"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
) // initialize client with credentials

func main() {
	workflow.New(&workflows.Asynchronous{}).Dispatch2()
}
