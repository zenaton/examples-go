package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var ParallelWorkflow = workflow.New("ParallelWorkflow",
	func() (interface{}, error) {

		//tasks A and B run in parallel
		err := task.Parallel{
			tasks.TaskA.NewInstance(),
			tasks.TaskB.NewInstance(),
		}.Execute()

		if err != nil {
			panic(err)
		}

		tasks.TaskC.NewInstance().Execute()
		return nil, nil
	})
