package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

func init() {
	//todo: call this registerworkflow, and check in newWorkflow that this is registered first
	workflow.New(&Event{})
}

type Event struct{}

func (w *Event) Handle() {
	task.New(&tasks.A{}).Execute()
	task.New(&tasks.B{}).Execute()
}

func (w *Event) OnEvent(eventName string, eventData interface{}) {
	if eventName == "MyEvent" {
		task.New(&tasks.C{}).Execute()
	}
}

func (w *Event) ID() string {
	return "MyId"
}
