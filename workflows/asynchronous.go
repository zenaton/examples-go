package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var AsynchronousWorkflow = workflow.New2(&Asynchronous{})

type Asynchronous struct{}

func (a Asynchronous) Handle() {
	tasks.TaskA().Dispatch()
	tasks.TaskB().Dispatch()
}
