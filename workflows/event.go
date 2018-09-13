package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var EventWorkflow = workflow.NewCustom("EventWorkflow", &Event{})

type Event struct {
	IDstr string
}

func (e *Event) Handle() (interface{}, error) {
	tasks.TaskA.NewInstance().Execute()
	tasks.TaskB.NewInstance().Execute()
	return nil, nil
}

func (e *Event) OnEvent(eventName string, eventData interface{}) {
	if eventName == "MyEvent" {
		tasks.TaskC.NewInstance().Execute()
	}
}

func (e *Event) ID() string {
	if e.IDstr == "" {
		return "MyId"
	}
	return e.IDstr
}
