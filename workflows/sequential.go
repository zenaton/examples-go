package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var SequentialWorkflow = workflow.New("SequentialWorkflow",
	func() (interface{}, error) {

		var a int
		tasks.A.New().Execute().Output(&a)

		if a > 0 {
			tasks.B.New().Execute()
		} else {
			tasks.C.New().Execute()
		}

		tasks.D.New().Execute()

		return nil, nil
	})
