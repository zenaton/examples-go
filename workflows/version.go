package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

//todo: I don't really understand what's happening here. Should figure it out so that I can document it.
//todo: this isn't right. Shouldn't have an instance of the version thing
var VersionWorkflow = workflow.NewVersion2("VersionWorkflow", []*workflow.Workflow{
	workflow.New(&v0{}),
	workflow.New(&v1{}),
	workflow.New(&v2{}),
})

type v0 struct{}

func (v v0) Handle() {
	task.Parallel{
		tasks.TaskA(),
		tasks.TaskB(),
	}.Execute()
}

type v1 struct{}

func (v v1) Handle() {
	task.Parallel{
		tasks.TaskA(),
		tasks.TaskB(),
		tasks.TaskC(),
	}.Execute()
}

type v2 struct{}

func (v v2) Handle() {
	task.Parallel{
		tasks.TaskA(),
		tasks.TaskB(),
		tasks.TaskC(),
		tasks.TaskD(),
	}.Execute()
}
