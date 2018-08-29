package workflow

import (
	"fmt"
	"github.com/zenaton/zenaton-go/v1/zenaton"
)

type recursiveData struct {
	ID  int
	Max int
}

var DisplayTask *zenaton.Task

var RelaunchTask *zenaton.Task

func init() {
	RelaunchTask = zenaton.NewTask(zenaton.TaskParams{
		Name: "RelaunchTask",
		HandleFunc: func(data recursiveData) {
			if data.ID >= data.Max {
				return
			}

			newID := 1 + data.ID
			fmt.Printf("\nIteration: %v\n", newID)

			RecursiveWorkflow.SetData(recursiveData{newID, data.Max}).Dispatch()
		},
	})

	DisplayTask = zenaton.NewTask(zenaton.TaskParams{
		Name: "DisplayTask",
		HandleFunc: func(data recursiveData) {
			fmt.Print(data.ID)
		},
	})
}

var RecursiveWorkflow = zenaton.NewWorkflow(zenaton.WorkflowParams{
	Name: "RecursiveWorkflow",
	HandleFunc: func(data recursiveData) {

		for counter := 0; counter < 1; counter++ {
			DisplayTask.SetData(recursiveData{ID: counter}).Execute()
		}

		RelaunchTask.SetData(data).Execute()
	},
})
