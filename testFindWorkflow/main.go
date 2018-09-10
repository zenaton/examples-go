package main

import (
	_ "github.com/zenaton/examples-go/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflows"
	"time"
)

func main() {

	workflows.TestW.NewInstance(workflows.Test2{
		IDstr:      "test",
		TaskReturn: "testTaskReturn",
	}).Dispatch()

	wf, err := workflows.TestW.WhereID("test").Find()

	if err != nil {
		panic(err)
	}

	instance := wf.GetData().(*workflows.Test2)

	//to make the logs more predictable between dispatches
	time.Sleep(5 * time.Second)

	// launch a new instance with same return value
	workflows.TestW.NewInstance(workflows.Test2{
		TaskReturn: instance.TaskReturn,
	}).Dispatch()
}
