package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var WaitEventWorkflow = workflow.NewCustom("WaitEventWorkflow", &WaitEvent{})

type WaitEvent struct{}

func (w *WaitEvent) Handle() (interface{}, error) {
	event, err := task.Wait().ForEvent("MyEvent").Seconds(4).Execute()

	if err != nil {
		panic(err)
	}

	if event.Arrived() {
		tasks.TaskA.NewInstance().Execute()
	} else {
		tasks.TaskB.NewInstance().Execute()
	}
	return nil, nil
}

func (w *WaitEvent) ID() string {
	return "MyId"
}
