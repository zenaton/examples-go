package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

func init() {
	//todo: call this registerworkflow, and check in newWorkflow that this is registered first
	workflow.New(&Wait{})
	task.Wait() //todo: god this is ugly
}

type Wait struct{}

func (w *Wait) Handle() {
	task.New(&tasks.A{}).Execute()
	task.Wait().ForEvent("MyEvent").Seconds(4).Execute()
	task.New(&tasks.D{}).Execute()
}
