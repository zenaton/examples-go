package workflows

import (
	"fmt"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

func init() {
	//todo: call this registerworkflow, and check in newWorkflow that this is registered first
	workflow.New(&Recursive{})
	task.New(&DisplayTask{})
	task.New(&RelaunchTask{})
}

type Recursive struct {
	ID  int
	Max int
}

func (r Recursive) Handle() {
	for counter := 0; counter < 1; counter++ {
		task.New(&DisplayTask{ID: counter}).Execute()
	}

	task.New(&RelaunchTask{R: r}).Execute()
}

type DisplayTask struct {
	ID int
}

func (dt *DisplayTask) Handle() {
	fmt.Print(dt.ID)
}

type RelaunchTask struct {
	R Recursive
}

func (rt *RelaunchTask) Handle() {
	if rt.R.ID >= rt.R.Max {
		return
	}

	newID := 1 + rt.R.ID
	fmt.Printf("\nIteration: %v\n", newID)

	workflow.New(&Recursive{newID, rt.R.Max}).Dispatch2()
}
