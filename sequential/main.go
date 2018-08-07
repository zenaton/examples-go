package main

import _ "github.com/zenaton/examples-go/client" // initialize zenaton client with credentials
import "github.com/zenaton/examples-go/workflow" // initialize zenaton client with credentials

func main() {
	workflow.SequentialWorkflow.Dispatch()
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
