package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

//set the default state to true
var EventWorkflow = workflow.NewCustom("EventWorkflow", &Event{State: true})

// All fields in this struct must be exported, as they must be serialized.
type Event struct {
	IDstr string
	State bool
}

func (e *Event) Handle() (interface{}, error) {
	tasks.TaskA.NewInstance().Execute()

	if e.State {
		tasks.TaskB.NewInstance().Execute()
	} else {
		tasks.TaskC.NewInstance().Execute()
	}

	return nil, nil
}

func (e *Event) OnEvent(eventName string, eventData interface{}) {

	if eventName == "MyEvent" {
		e.State = false
	}
}

func (e *Event) ID() string {
	return e.IDstr
}
