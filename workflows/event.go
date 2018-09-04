package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var EventWorkflow = workflow.New(&Event{})

type Event struct{}

func (w *Event) Handle() (interface{}, error) {
	tasks.TaskA.NewInstance().Execute()
	tasks.TaskB.NewInstance().Execute()
	return nil, nil
}

func (w *Event) OnEvent(eventName string, eventData interface{}) {
	if eventName == "MyEvent" {
		tasks.TaskC.NewInstance().Execute()
	}
}

func (w *Event) ID() string {
	return "MyId"
}
