package workflow

import (
	"github.com/zenaton/examples-go/task"
	"github.com/zenaton/zenaton-go/v1/zenaton"
)

var WaitWorkflow = zenaton.NewWorkflow(zenaton.WorkflowParams{
	Name: "WaitWorkflow",
	HandleFunc: func() {
		// todo: figure out how to do something like this.email in javascript example
		task.A.Execute()
		// todo: kind of ugly to pass in nil here, maybe do a .data?
		zenaton.NewWait().Seconds(5).Execute()
		task.B.Execute()
	},
})
