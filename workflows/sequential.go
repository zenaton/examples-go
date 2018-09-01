package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var SequentialWorkflow = workflow.New2(&Sequential{})

type Sequential struct{}

func (s Sequential) Handle() {
	tasks.TaskA().Execute()
	tasks.TaskB().Execute()
}
