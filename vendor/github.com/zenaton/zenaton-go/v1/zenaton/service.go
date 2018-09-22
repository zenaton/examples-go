package zenaton

import (
	"github.com/zenaton/zenaton-go/v1/zenaton/client"
	"github.com/zenaton/zenaton-go/v1/zenaton/engine"
	"github.com/zenaton/zenaton-go/v1/zenaton/errors"
	"github.com/zenaton/zenaton-go/v1/zenaton/interfaces"
	"github.com/zenaton/zenaton-go/v1/zenaton/service/serializer"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

type Service struct {
	Serializer      *serializer.Serializer
	Client          *client.Client
	Engine          *engine.Engine
	WorkflowManager *workflow.Store
	TaskManager     *task.Store
	Errors          Errors
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

func NewService() *Service {
	return &Service{
		Client:          client.NewClient(true),
		Engine:          engine.NewEngine(),
		Serializer:      &serializer.Serializer{},
		WorkflowManager: workflow.Manager,
		TaskManager:     task.Manager,
		Errors: Errors{
			ScheduledBoxError:    errors.ScheduledBoxError,
			ExternalZenatonError: errors.ExternalZenatonError,
			InternalZenatonError: errors.InternalZenatonError,
		},
	}
}

func InitClient(appID, apiToken, appEnv string) {
	client.InitClient(appID, apiToken, appEnv)
}

type Workflow = workflow.Instance
type Task = task.Instance
type Wait = task.WaitTask
type Job = interfaces.Job
type Processor = engine.Processor
