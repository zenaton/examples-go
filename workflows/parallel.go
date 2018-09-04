package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var ParallelWorkflow = workflow.New(&Parallel{})

type Parallel struct{}

func (s Parallel) Handle() (interface{}, error) {

	//tasks A and B run in parallel
	_, err := task.Parallel{
		tasks.TaskA.NewInstance(),
		tasks.TaskB.NewInstance(),
	}.Execute()

	if err != nil {
		panic(err)
	}

	tasks.TaskC.NewInstance().Execute()
	return nil, nil
}
