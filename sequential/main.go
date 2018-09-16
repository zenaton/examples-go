package main

import (
	_ "github.com/zenaton/examples-go" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflows"
)

func main() {

	workflows.SequentialWorkflow.New().Dispatch()
}
