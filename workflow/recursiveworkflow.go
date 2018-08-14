package workflow

import (
	"strconv"

	"github.com/zenaton/examples-go/idmax"
	"github.com/zenaton/zenaton-go/v1/zenaton"
)

var RecursiveWorkflow = zenaton.NewWorkflow(zenaton.WorkflowParams{
	Name: "RecursiveWorkflow",
	HandleFunc: func(i idmax.IDmax) {

		for counter := 0; counter < 3; counter++ {
			DisplayTask.SetData(idmax.IDmax{ID: counter}).Execute()
		}

		RelaunchTask.SetData(i).Execute()
	},
	ID: func(i idmax.IDmax) string {
		return strconv.Itoa(i.ID)
	},
})
