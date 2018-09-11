package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var VersionWorkflow = workflow.Version("VersionWorkflow", []*workflow.WorkflowType{
	V0Workflow,
	V1Workflow,
	V2Workflow,
})

var V0Workflow = workflow.NewDefault("V0", func() (interface{}, error) {
	task.Parallel{
		tasks.TaskA.NewInstance(),
		tasks.TaskB.NewInstance(),
	}.Execute()
	return nil, nil
})

var V1Workflow = workflow.NewDefault("V1", func() (interface{}, error) {
	task.Parallel{
		tasks.TaskA.NewInstance(),
		tasks.TaskB.NewInstance(),
		tasks.TaskC.NewInstance(),
	}.Execute()
	return nil, nil
})

var V2Workflow = workflow.NewDefault("V2", func() (interface{}, error) {
	task.Parallel{
		tasks.TaskA.NewInstance(),
		tasks.TaskB.NewInstance(),
		tasks.TaskC.NewInstance(),
		tasks.TaskD.NewInstance(),
	}.Execute()
	return nil, nil
})
