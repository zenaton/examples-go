package main

import (
	"github.com/zenaton/examples-go/integration/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflows"
)

func init() { client.SetEnv("version.env") }
func main() {
	workflows.VersionWorkflow.NewInstance().Dispatch()
}
