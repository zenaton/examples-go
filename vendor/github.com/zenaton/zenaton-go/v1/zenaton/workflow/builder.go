package workflow

import (
	"github.com/zenaton/zenaton-go/v1/zenaton/internal/client"
)

// The QueryBuilder allows you to Find, Kill, Pause, and Resume workflow instances by id. You can also Send an event
// to a workflow.
type QueryBuilder struct {
	workflowDefinition string
	id                 string
	client             *client.Client
}

func newBuilder(name string) *QueryBuilder {
	return &QueryBuilder{
		client:             client.NewClient(false),
		workflowDefinition: name,
	}
}

func (b *QueryBuilder) whereID(id string) *QueryBuilder {
	b.id = id
	return b
}

// Find allows you to find a running instance of a workflow. If no instance with the provided id (from WhereID) is found
// Find will return nil, nil. You will only get a non-nil error if there is a problem with the http request sent
// to retrieve the instance.
func (b *QueryBuilder) Find() (*Instance, error) {
	output, ok, err := b.client.FindWorkflowInstance(b.workflowDefinition, b.id)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	properties := output["data"]["properties"]
	name := output["data"]["name"]

	return UnsafeManager.UnsafeGetInstance(name, properties)
}

// Send an event to a workflow.
func (b *QueryBuilder) Send(eventName string, eventData interface{}) {
	b.client.SendEvent(b.workflowDefinition, b.id, eventName, eventData)
}

// Kill a workflowDef instance
func (b *QueryBuilder) Kill() (*QueryBuilder, error) {
	err := b.client.KillWorkflow(b.workflowDefinition, b.id)
	return b, err
}

// Pause a workflowDef instance
func (b *QueryBuilder) Pause() (*QueryBuilder, error) {
	err := b.client.PauseWorkflow(b.workflowDefinition, b.id)
	return b, err
}

// Resume a workflowDef instance
func (b *QueryBuilder) Resume() (*QueryBuilder, error) {
	err := b.client.ResumeWorkflow(b.workflowDefinition, b.id)
	return b, err
}
