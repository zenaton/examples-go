package workflow

import (
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var RecursiveWorkflow = workflow.NewCustom("RecursiveWorkflow", &Recursive{})

type Recursive struct {
	ID  int
	Max int
}

func (r *Recursive) Init(id, max int) {
	r.ID = id
	r.Max = max
}

func (r *Recursive) Handle() (interface{}, error) {

	for counter := 0; counter < 2; counter++ {
		DisplayTask.New(counter).Execute()
	}

	RelaunchTask.New(r).Execute()
	return nil, nil
}
