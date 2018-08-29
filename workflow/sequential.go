package workflow

import (
	"github.com/zenaton/examples-go/task"
	"github.com/zenaton/zenaton-go/v1/zenaton"
)

var SequentialWorkflow = zenaton.NewWorkflow(zenaton.WorkflowParams{
	Name: "SequentialWorkflow",
	HandleFunc: func() {
		task.A.Execute()
		task.B.Execute()
	},
})
