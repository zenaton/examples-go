package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var WaitEventWorkflow = workflow.NewCustom("WaitEventWorkflow", &WaitEvent{})

type WaitEvent struct {
	IDstr string
}

func (w *WaitEvent) Handle() (interface{}, error) {

	// Waits until the event or at most 4 seconds
	event := task.Wait().ForEvent("MyEvent").Seconds(4).Execute()

	if event != nil {
		// if event has been triggered within 4 seconds
		tasks.TaskA.NewInstance().Execute()
	} else {
		//else
		tasks.TaskB.NewInstance().Execute()
	}
	return nil, nil
}

func (w *WaitEvent) ID() string {
	return w.IDstr
}
