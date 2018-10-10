package zenaton

import (
	"github.com/zenaton/zenaton-go/v1/zenaton/errors"
	"github.com/zenaton/zenaton-go/v1/zenaton/internal/client"
	"github.com/zenaton/zenaton-go/v1/zenaton/internal/engine"
	"github.com/zenaton/zenaton-go/v1/zenaton/service/serializer"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

// UnsafeService contains many things that the agent needs to operate.
// For example, when defining a new workflow Definition, the definition will be stored in the WorkflowManager.
// The agent will use the WorkflowManager to retrieve a workflow Definition given a workflow name.
type UnsafeService struct {
	// None of these Values should be used directly and can lead to unpredicted behaviors.
	Serializer      *serializer.Serializer
	Client          *client.Client
	Engine          *engine.Engine
	WorkflowManager *workflow.Store
	TaskManager     *task.Store
	Errors          Errors
}

// NewService creates a new Zenaton service.
// In your boot file, you must have this line (exactly): "var Service = zenaton.NewService()"
func NewService() *UnsafeService {
	return &UnsafeService{
		Client:          client.NewClient(true),
		Engine:          engine.NewEngine(),
		Serializer:      &serializer.Serializer{},
		WorkflowManager: workflow.UnsafeManager,
		TaskManager:     task.UnsafeManager,
		Errors: Errors{
			ScheduledBoxError:    errors.ScheduledBoxError,
			ExternalZenatonError: errors.ExternalZenatonError,
			InternalZenatonError: errors.InternalZenatonError,
		},
	}
}

// InitClient will initialize the Zenaton client with your credentials and app env.
func InitClient(appID, apiToken, appEnv string) {
	client.InitClient(appID, apiToken, appEnv)
}

// Errors is provided so that the agent can use the Errors package without the user of the library having to re-export
// the errors package
type Errors struct {
	ScheduledBoxError    string
	ExternalZenatonError string
	InternalZenatonError string
}

func (e *Errors) New(name, message string) errors.ZenatonError {
	return errors.NewWithOffset(name, message, 4)
}

func (e *Errors) Wrap(name string, err error) errors.ZenatonError {
	return errors.WrapWithOffset(name, err, 4)
}

// Workflow aliases a workflow.Instance
type Workflow = workflow.Instance

// Task aliases a task.Instance
type Task = task.Instance

// Wait aliases a task.WaitTask
type Wait = task.WaitTask

// Job aliases a engine.Job
type Job = engine.Job

// Processor aliases a engine.Processor
type Processor = engine.Processor
