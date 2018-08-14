package workflow

import (
	"fmt"
	"strconv"
	"time"

	"github.com/zenaton/examples-go/idmax"
	"github.com/zenaton/zenaton-go/v1/zenaton"
)

var DisplayTask *zenaton.Task

var RelaunchTask *zenaton.Task

func init() {
	RelaunchTask = zenaton.NewTask(zenaton.TaskParams{
		Name: "RelaunchTask",
		HandleFunc: func(idMax idmax.IDmax) {
			if idMax.ID >= idMax.Max {
				return
			}

			newID := 1 + idMax.ID
			fmt.Printf("\nIteration: %v\n", newID)

			RecursiveWorkflow.SetData(idmax.IDmax{newID, idMax.Max}).Dispatch()
		},

		//todo: shouldn't need to do strconv here should I?
		//todo: should have this be like the handle func where you can define any input
		ID: func(idMax interface{}) string {
			return strconv.Itoa(idMax.(idmax.IDmax).ID)
		},
	})

	DisplayTask = zenaton.NewTask(zenaton.TaskParams{
		Name: "DisplayTask",
		HandleFunc: func(i idmax.IDmax) {
			fmt.Print(i.ID)
			time.Sleep(1 * time.Second)
		},
		ID: func(idMax interface{}) string {
			return string(idMax.(idmax.IDmax).ID)
		},
	})
}
