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
	tasks.A.New().Execute()

	if e.State {
		tasks.B.New().Execute()
	} else {
		tasks.C.New().Execute()
	}

	return nil, nil
}

func (e *Event) Init(id string) {
	e.IDstr = id
}

func (e *Event) OnEvent(eventName string, eventData interface{}) {

	if eventName == "MyEvent" {
		e.State = false
	}
}

func (e *Event) ID() string {
	return e.IDstr
}
