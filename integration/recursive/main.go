package main

import (
	// initialize zenaton client with credentials

	"github.com/zenaton/examples-go/integration/client"
	"github.com/zenaton/examples-go/workflows"
)

func init() { client.SetEnv("recursive.env") }
func main() {
	workflows.RecursiveWorkflow.NewInstance(workflows.Recursive{0, 2}).Dispatch()
}
