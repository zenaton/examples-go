package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var SequentialWorkflow = workflow.NewDefault("SequentialWorkflow",
	func() (interface{}, error) {
		tasks.TaskA.NewInstance().Execute()
		tasks.TaskB.NewInstance().Execute()
		return nil, nil
	})
