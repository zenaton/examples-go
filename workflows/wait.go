package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var WaitWorkflow = workflow.New("WaitWorkflow", &Wait{})

type Wait struct{}

func (w *Wait) Handle() (interface{}, error) {
	tasks.TaskA.NewInstance().Execute()
	task.Wait().ForEvent("MyEvent").Seconds(4).Execute()
	tasks.TaskD.NewInstance().Execute()
	return nil, nil
}
