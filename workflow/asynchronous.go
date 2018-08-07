package workflow

import (
	"github.com/zenaton/examples-go/task"
	"github.com/zenaton/zenaton-go/v1/zenaton"
)

var AsynchronousWorkflow = zenaton.NewWorkflow(zenaton.WorkflowParams{
	Name: "AsynchronousWorkflow",
	HandleFunc: func() {
		task.A.Dispatch()
		task.B.Execute()
	},
})
