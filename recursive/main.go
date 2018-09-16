package main

import (
	// initialize zenaton client with credentials

	_ "github.com/zenaton/examples-go"
	"github.com/zenaton/examples-go/recursive/workflow"
)

func main() {
	workflow.RecursiveWorkflow.New(0, 2).Dispatch()
}
