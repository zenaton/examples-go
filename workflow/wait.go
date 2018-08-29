package workflow

import (
	"github.com/zenaton/examples-go/task"
	"github.com/zenaton/zenaton-go/v1/zenaton"
)

var WaitWorkflow = zenaton.NewWorkflow(zenaton.WorkflowParams{
	Name: "WaitWorkflow",
	HandleFunc: func() {
		task.A.Execute()
		zenaton.NewWait().Seconds(4).Execute()
		task.D.Execute()
	},
})
