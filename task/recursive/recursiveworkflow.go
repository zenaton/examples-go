package recursive

import (
	"strconv"

	"github.com/zenaton/zenaton-go/v1/zenaton"
)

func NewRecursiveWorkflow(id, max int) *zenaton.Workflow {

	type IDmax struct {
		ID  int
		Max int
	}

	handlerFunc := func(idmax IDmax) {

		for counter := 0; counter < 3; counter++ {
			NewDisplayTask(strconv.Itoa(counter)).Execute()
		}

		NewRelaunchTask(idmax.ID, idmax.Max).Execute()
	}

	idFunc := func() string {
		return strconv.Itoa(id)
	}

	idmax := IDmax{id, max}

	wf := zenaton.NewWorkflow("RecursiveWorkflow", handlerFunc).IDFunc(idFunc).Data(idmax)

	return wf
}
