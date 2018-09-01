package main

import (
	// initialize zenaton client with credentials

	_ "github.com/zenaton/examples-go/client"
	"github.com/zenaton/examples-go/workflows"
)

func main() {
	workflows.RecursiveWorkflow(workflows.Recursive{0, 2}).Dispatch()
}
