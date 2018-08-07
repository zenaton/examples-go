package workflow

//

//todo: should have the equivilent documentation of this from cadence:
/**
* This sample workflow demonstrates how to use multiple Cadence corotinues (instead of native goroutine) to process a
* chunk of a large work item in parallel, and then merge the intermediate result to generate the final result.
* In cadence workflow, you should not use go routine. Instead, you use corotinue via workflow.Go method.
 */

import (
	"github.com/zenaton/examples-go/task"
	"github.com/zenaton/zenaton-go/v1/zenaton"
)

var ParallelWorkflow = zenaton.NewWorkflow(zenaton.WorkflowParams{
	Name: "ParallelWorkflow",
	HandleFunc: func() {
		// tasks A and B run in parallel
		_, err := zenaton.Tasks{task.A, task.B}.Execute()
		if err != nil {
			panic(err)
		}

		task.C.Execute()
	},
})
