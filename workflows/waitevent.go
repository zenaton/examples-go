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

func (w *WaitEvent) Init(id string) {
	w.IDstr = id
}

func (w *WaitEvent) Handle() (interface{}, error) {

	execution := task.Wait().ForEvent("MyEvent").Seconds(4).Execute()

	if execution.EventReceived() {
		tasks.A.New().Execute()
	} else {
		tasks.B.New().Execute()
	}
	return nil, nil
}

func (w *WaitEvent) ID() string {
	return w.IDstr
}
