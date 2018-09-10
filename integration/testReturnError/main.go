package main

import (
	"github.com/zenaton/examples-go/integration/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflows"
)

func init() { client.SetEnv("testReturnError.env") }
func main() {

	workflows.TestW.WhereID("MyID").Kill()
	workflows.TestW.NewInstance(workflows.Test2{
		TaskError: "testTaskError",
	}).Dispatch()
}
