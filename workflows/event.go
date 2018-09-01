package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var EventWorkflow = workflow.New2(&Event{})

type Event struct{}

func (w *Event) Handle() {
	tasks.TaskA().Execute()
	tasks.TaskB().Execute()
}

func (w *Event) OnEvent(eventName string, eventData interface{}) {
	if eventName == "MyEvent" {
		tasks.TaskC().Execute()
	}
}

func (w *Event) ID() string {
	return "MyId"
}
