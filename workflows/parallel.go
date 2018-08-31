package workflows

//

//todo: should have the equivilent documentation of this from cadence:
/**
* This sample workflow demonstrates how to use multiple Cadence corotinues (instead of native goroutine) to process a
* chunk of a large work item in parallel, and then merge the intermediate result to generate the final result.
* In cadence workflow, you should not use go routine. Instead, you use corotinue via workflow.Go method.
 */

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

func init() {
	//todo: call this registerworkflow, and check in newWorkflow that this is registered first
	workflow.New(&Parallel{})
}

type Parallel struct{}

func (s Parallel) Handle() {
	//tasks A and B run in parallel
	_, err := task.Parallel{
		task.New(&tasks.A{}),
		task.New(&tasks.B{}),
	}.Execute()

	//todo: should I have errors?
	if err != nil {
		panic(err)
	}

	task.New(&tasks.C{}).Execute()
}
