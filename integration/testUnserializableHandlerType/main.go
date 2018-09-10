package main

import (
	"github.com/zenaton/examples-go/integration/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/examples-go/workflows"
)

func init() { client.SetEnv("testUnserializableHandlerType.env") }
func main() {
	//workflows.TestW.WhereID("MyID").Kill()
	t := tasks.TestTask.NewInstance()

	workflows.WithTaskWorkflow.NewInstance(&workflows.WithTask{t}).Dispatch()
}
