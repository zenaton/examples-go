/*
the boot file must have two things:
	1) an ignored import of all of the workflows. all workflows you wish to use must be exported, package level
		variables, initialized with workflow.New() or workflow.NewCustom(). See for example in: github.com/zenaton/examples-go/workflows/workflows
	2) an initialization of the zenaton client.
*/

package boot

import (
	// (1)
	// if this project is in your GOPATH, make sure to have imports that start from your GOPATH.
	// if this project is not in your GOPATH, then you can use relative imports here
	_ "github.com/zenaton/examples-go/workflows"
	// (2)
	_ "github.com/zenaton/examples-go/client"
	"github.com/zenaton/zenaton-go/v1/zenaton"
)

type Workflow = zenaton.Workflow
type Task = zenaton.Task
type Wait = zenaton.Wait
type Job = zenaton.Job
type Processor = zenaton.Processor

var Service = zenaton.NewService()
