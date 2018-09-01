package workflows

import (
	"fmt"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var RecursiveWorkflow = workflow.New2(&Recursive{})

type Recursive struct {
	ID  int
	Max int
}

func (r Recursive) Handle() {
	for counter := 0; counter < 1; counter++ {
		DisplayTask(&Display{ID: counter}).Execute()
	}

	RelaunchTask(&Relaunch{R: r}).Execute()
}

var DisplayTask = task.New2(&Display{})

type Display struct {
	ID int
}

func (dt *Display) Handle() {
	fmt.Print(dt.ID)
}

var RelaunchTask = task.New2(&Relaunch{})

type Relaunch struct {
	R Recursive
}

func (rt *Relaunch) Handle() {
	if rt.R.ID >= rt.R.Max {
		return
	}

	newID := 1 + rt.R.ID
	fmt.Printf("\nIteration: %v\n", newID)

	workflow.New(&Recursive{newID, rt.R.Max}).Dispatch()
}
