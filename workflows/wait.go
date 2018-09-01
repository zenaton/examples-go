package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var WaitWorkflow = workflow.New2(&Wait{})

type Wait struct{}

func (w *Wait) Handle() {
	tasks.TaskA().Execute()
	task.Wait().ForEvent("MyEvent").Seconds(4).Execute()
	tasks.TaskD().Execute()
}
