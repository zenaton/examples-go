package main

import (
	// initialize zenaton client with credentials
	_ "github.com/zenaton/examples-go/client"
	"github.com/zenaton/examples-go/idmax"
	"github.com/zenaton/examples-go/workflow"
)

func main() {
	workflow.RecursiveWorkflow.SetData(idmax.IDmax{0, 2}).Dispatch()
}
