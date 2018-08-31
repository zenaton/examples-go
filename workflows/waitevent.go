package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

//var WaitEventWorkflow = workflow.New(zenaton.WorkflowParams{
//	Name: "WaitEventWorkflow",
//	HandleFunc: func() {
//		//todo:
//		event := zenaton.NewWait().Seconds(4).Execute()
//		task.A.Execute()
//		zenaton.NewWait().Timestamp(time.Now().Unix() + 5).Execute()
//		task.B.Execute()
//	},
//})

func init() {
	//todo: call this registerworkflow, and check in newWorkflow that this is registered first
	//todo: test that you get a propper error message if you ddon't do this
	workflow.New(&WaitEvent{})
	task.Wait() //todo: get rid of this
}

type WaitEvent struct{}

func (w *WaitEvent) Handle() {
	event, _ := task.Wait().ForEvent("MyEvent").Seconds(4).Execute()

	if event != nil {
		task.New(&tasks.A{}).Execute()
	} else {
		task.New(&tasks.B{}).Execute()
	}
}

func (w *WaitEvent) ID() string {
	return "MyId"
}
