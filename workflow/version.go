package workflow

import (
	"github.com/zenaton/examples-go/task"
	"github.com/zenaton/zenaton-go/v1/zenaton"
)

//todo: I don't really understand what's happening here. Should figure it out so that I can document it.
var VersionWorkflow = zenaton.NewVersion("VersionWorkflow", []*zenaton.Workflow{v0, v1, v2})

var v0 = zenaton.NewWorkflow(zenaton.WorkflowParams{
	Name: "VersionWorkflow_v0",
	HandleFunc: func() {
		zenaton.Tasks{
			task.A,
			task.B,
		}.Execute()
	},
})

var v1 = zenaton.NewWorkflow(zenaton.WorkflowParams{
	Name: "VersionWorkflow_v1",
	HandleFunc: func() {
		zenaton.Tasks{
			task.A,
			task.B,
			task.C,
		}.Execute()
	},
})

var v2 = zenaton.NewWorkflow(zenaton.WorkflowParams{
	Name: "VersionWorkflow_v2",
	HandleFunc: func() {
		zenaton.Tasks{
			task.A,
			task.B,
			task.C,
			task.D,
		}.Execute()
	},
})
