package workflow

import (
	"fmt"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var RelaunchTask = task.NewCustom("RelaunchTask", &Relaunch{})

type Relaunch struct {
	R *Recursive
}

func (rt *Relaunch) Init(recursive *Recursive) {
	rt.R = recursive
}

func (rt *Relaunch) Handle() (interface{}, error) {

	if rt.R.ID >= rt.R.Max {
		return nil, nil
	}

	newID := 1 + rt.R.ID
	fmt.Println("\nIteration:", newID)

	RecursiveWorkflow.New(newID, rt.R.Max).Dispatch()
	return nil, nil
}
