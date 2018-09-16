package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var VersionWorkflow = workflow.Version("VersionWorkflow", []*workflow.Definition{
	V0Workflow,
	V1Workflow,
	V2Workflow,
})

var V0Workflow = workflow.New("V0", func() (interface{}, error) {
	task.Parallel{
		tasks.A.New(),
		tasks.B.New(),
	}.Execute()
	return nil, nil
})

var V1Workflow = workflow.New("V1", func() (interface{}, error) {
	task.Parallel{
		tasks.A.New(),
		tasks.B.New(),
		tasks.C.New(),
	}.Execute()
	return nil, nil
})

var V2Workflow = workflow.New("V2", func() (interface{}, error) {
	task.Parallel{
		tasks.A.New(),
		tasks.B.New(),
		tasks.C.New(),
		tasks.D.New(),
	}.Execute()
	return nil, nil
})
