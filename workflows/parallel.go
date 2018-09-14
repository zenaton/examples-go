package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var ParallelWorkflow = workflow.New("ParallelWorkflow",
	func() (interface{}, error) {

		var a int
		var b int

		//tasks A and B run in parallel
		task.Parallel{
			tasks.TaskA.NewInstance(),
			tasks.TaskB.NewInstance(),
		}.Execute().Output(&a, &b)

		if a > b {
			tasks.TaskC.NewInstance().Execute()
		} else {
			tasks.TaskD.NewInstance().Execute()
		}

		return nil, nil
	})
