package main

import (
	_ "github.com/zenaton/examples-go/client"
	"github.com/zenaton/examples-go/workflows"
) // initialize client with credentials

func main() {
	workflows.AsynchronousWorkflow().Dispatch()
}
