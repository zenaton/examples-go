package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var SequentialWorkflow = workflow.New("SequentialWorkflow", &Sequential{})

type Sequential struct{}

func (s Sequential) Handle() (interface{}, error) {
	tasks.TaskA.NewInstance().Execute()
	tasks.TaskB.NewInstance().Execute()
	return nil, nil
}
