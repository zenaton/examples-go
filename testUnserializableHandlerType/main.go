package main

import (
	_ "github.com/zenaton/examples-go/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/examples-go/workflows"
)

func main() {
	//workflows.TestW.WhereID("MyID").Kill()
	t := tasks.TestTask.NewInstance()

	workflows.WithTaskWorkflow.NewInstance(&workflows.WithTask{t}).Dispatch()
}
