package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var SequentialWorkflow = workflow.New("SequentialWorkflow",
	func() (interface{}, error) {

		var a int
		tasks.TaskA.NewInstance().Execute(&a)

		if a == 0 {
			tasks.TaskB.NewInstance().Execute()
		} else {
			tasks.TaskC.NewInstance().Execute()
		}

		return nil, nil
	})
