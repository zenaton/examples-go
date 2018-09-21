/*
the boot file must have two things:
	1) an ignored import of all of the workflows. All workflows you wish to use must be exported, package level
		variables, initialized with workflow.New() or workflow.NewCustom(). See for example in: github.com/zenaton/examples-go/workflows/workflows
	2) an initialization of the zenaton client.
	3) re-export all the necessary types and information from the zenaton library so that the agent can use them.
*/

package boot

import (
	// (1)
	_ "github.com/zenaton/examples-go/recursive/workflow"
	_ "github.com/zenaton/examples-go/workflows"
	// (2)
	_ "github.com/zenaton/examples-go"
	"github.com/zenaton/zenaton-go/v1/zenaton"
)

// (3) these lines must all be present in your boot file.
type Workflow = zenaton.Workflow
type Task = zenaton.Task
type Wait = zenaton.Wait
type Job = zenaton.Job
type Processor = zenaton.Processor

var Service = zenaton.NewService()
