package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var MaxProcessingWorkflow = workflow.NewDefault("MaxProcessingWorkflow",
	func() (interface{}, error) {

		//tasks A and B run in parallel
		tasks.MaxProcessingTask.NewInstance(&tasks.MaxProcessing{
			Max: 2,
		}).Execute()
		return nil, nil
	})
