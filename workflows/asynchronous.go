package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var AsynchronousWorkflow = workflow.New("AsynchronousWorkflow", &Asynchronous{})

type Asynchronous struct{}

func (a Asynchronous) Handle() (interface{}, error) {
	tasks.TaskA.NewInstance().Dispatch()
	tasks.TaskB.NewInstance().Dispatch()
	return nil, nil
}
