package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var ErrorWorkflow = workflow.New("ErrorWorkflow",
	func() (interface{}, error) {

		//tasks A and E run in parallel
		task.Parallel{
			tasks.A.New(),
			//tasks.E panics
			tasks.E.New(),
		}.Execute()

		tasks.C.New().Execute()

		return nil, nil
	})
