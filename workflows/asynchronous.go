package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var AsynchronousWorkflow = workflow.NewDefault("AsynchronousWorkflow",
	func() (interface{}, error) {
		tasks.TaskA.NewInstance().Dispatch()
		tasks.TaskB.NewInstance().Dispatch()
		return nil, nil
	})
