package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var AsynchronousWorkflow = workflow.New("AsynchronousWorkflow",
	func() (interface{}, error) {

		tasks.A.New().Dispatch()
		tasks.B.New().Dispatch()

		tasks.C.New().Execute()
		tasks.D.New().Execute()

		return nil, nil
	})
