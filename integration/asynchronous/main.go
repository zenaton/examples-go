package main

import (
	"github.com/zenaton/examples-go/integration/client"
	"github.com/zenaton/examples-go/workflows"
) // initialize client with credentials
func init(){client.SetEnv("asynchronous.env")}
func main() {
	workflows.AsynchronousWorkflow.NewInstance().Dispatch()
}
