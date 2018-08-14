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

//todo: change parallel to work synchronously per conversation with Gilles
//todo: make sure there are no race conditions in the case of running these things concurrently

////workflows.AsynchronousWorkflow.Dispatch()
////workflows.ParallelWorkflow.Dispatch()
//
//workflows.EventWorkflow.Dispatch()
//time.Sleep(2 * time.Second)
//workflows.EventWorkflow.WhereID("MyId").Send("MyEvent", nil)
//
////workflows.WaitWorkflow.Dispatch()
//
////workflows.WaitEventWorkflow.Dispatch()
////time.Sleep(2 * time.Second)
////workflows.WaitEventWorkflow.WhereID("MyId").Send("MyEvent", nil)
//
////recursive.NewRecursiveWorkflow(0, 2).Dispatch()
////workflows.VersionWorkflow.Dispatch()
//
//time.Sleep(4 * time.Second)
