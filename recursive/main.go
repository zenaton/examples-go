package main

import (
	// initialize zenaton client with credentials

	_ "github.com/zenaton/examples-go/client"
	"github.com/zenaton/examples-go/workflows"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

func main() {
	workflow.New(&workflows.Recursive{0, 2}).Dispatch2()
}
