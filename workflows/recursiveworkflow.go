package workflows

import (
	"fmt"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var RecursiveWorkflow = workflow.New(&Recursive{})

type Recursive struct {
	ID  int
	Max int
}

func (r Recursive) Handle() (interface{}, error) {
	for counter := 0; counter < 1; counter++ {
		DisplayTask.NewInstance(&Display{ID: counter}).Execute()
	}

	RelaunchTask.NewInstance(&Relaunch{R: r}).Execute()
	return nil, nil
}

var DisplayTask = task.New(&Display{})

type Display struct {
	ID int
}

func (dt *Display) Handle() (interface{}, error) {
	fmt.Print(dt.ID)
	return nil, nil
}

var RelaunchTask = task.New(&Relaunch{})

type Relaunch struct {
	R Recursive
}

func (rt *Relaunch) Handle() (interface{}, error) {
	if rt.R.ID >= rt.R.Max {
		return nil, nil
	}

	newID := 1 + rt.R.ID
	fmt.Printf("\nIteration: %v\n", newID)

	RecursiveWorkflow.NewInstance(&Recursive{newID, rt.R.Max}).Dispatch()
	return nil, nil
}
