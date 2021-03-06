package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var WaitWorkflow = workflow.New("WaitWorkflow",
	func() (interface{}, error) {

		tasks.A.New().Execute()

		task.Wait().Seconds(5).Execute()

		tasks.B.New().Execute()

		return nil, nil
	})
