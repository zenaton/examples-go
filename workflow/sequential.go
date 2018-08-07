package workflow

import (
	"github.com/zenaton/examples-go/task"
	"github.com/zenaton/zenaton-go/v1/zenaton"
)

var SequentialWorkflow = zenaton.NewWorkflow(zenaton.WorkflowParams{
	Name: "SequentialWorkflow",
	HandleFunc: func() {
		//todo: should I use a pointer like in json.unmarshal? That way I wouldn't have to do a type assertion on output
		//_, err := task.A.Execute()
		//if err != nil {
		//	fmt.Println("error in sequential workflow: ", err)
		//}
		task.A.Execute()
		task.B.Execute()
	},
})
