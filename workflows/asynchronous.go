package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

func init() {
	workflow.New(&Asynchronous{})
}

type Asynchronous struct{}

func (a Asynchronous) Handle() {
	task.New(&tasks.A{}).Dispatch()
	task.New(&tasks.B{}).Dispatch()
}
