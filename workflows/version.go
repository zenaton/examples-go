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

var V0Workflow = workflow.New("V0", &V0{})

type V0 struct{}

func (v V0) Handle() (interface{}, error) {
	task.Parallel{
		tasks.TaskA.NewInstance(),
		tasks.TaskB.NewInstance(),
	}.Execute()
	return nil, nil
}

var V1Workflow = workflow.New("V1", &V1{})

type V1 struct{}

func (v V1) Handle() (interface{}, error) {
	task.Parallel{
		tasks.TaskA.NewInstance(),
		tasks.TaskB.NewInstance(),
		tasks.TaskC.NewInstance(),
	}.Execute()
	return nil, nil
}

var V2Workflow = workflow.New("V2", &V2{})

type V2 struct{}

func (v V2) Handle() (interface{}, error) {
	task.Parallel{
		tasks.TaskA.NewInstance(),
		tasks.TaskB.NewInstance(),
		tasks.TaskC.NewInstance(),
		tasks.TaskD.NewInstance(),
	}.Execute()
	return nil, nil
}
