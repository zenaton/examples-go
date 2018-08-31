package workflows

import (
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

func init() {
	workflow.New(&Sequential{})
}

type Sequential struct {
	Field string
}

func (s Sequential) Handle() {
	task.New(&tasks.A{}).Execute()
	task.New(&tasks.B{}).Execute()
}
