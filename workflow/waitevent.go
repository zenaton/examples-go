package workflow

import "github.com/zenaton/zenaton-go/v1/zenaton"

var WaitEventWorkflow = zenaton.NewWorkflow(zenaton.WorkflowParams{
	Name: "WaitEventWorkflow",
	HandleFunc: func() {
		//todo:
		//event := zenaton.NewWait().Seconds(4).Execute()
		//task.A.Execute()
		//zenaton.NewWait().Timestamp(time.Now().Unix() + 5).Execute()
		//task.B.Execute()
	},
})
