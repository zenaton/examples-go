package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
	"fmt"
)

//todo: test that you get a propper error message if you don't do this
var WaitEventWorkflow = workflow.New2(&WaitEvent{})

type WaitEvent struct{}

func (w *WaitEvent) Handle() {
	event, _ := task.Wait().ForEvent("MyEvent").Seconds(4).Execute()

	fmt.Println("event: ", event)

	if event != nil {
		tasks.TaskA().Execute()
	} else {
		tasks.TaskB().Execute()
	}
}

func (w *WaitEvent) ID() string {
	return "MyId"
}
