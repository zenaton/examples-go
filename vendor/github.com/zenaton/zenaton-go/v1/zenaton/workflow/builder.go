package workflow

import (
	"github.com/zenaton/zenaton-go/v1/zenaton/client"
)

type queryBuilder struct {
	workflowDefinition string
	id                 string
	client             *client.Client
}

func newBuilder(name string) *queryBuilder {
	return &queryBuilder{
		client:             client.NewClient(false),
		workflowDefinition: name,
	}
}

func (b *queryBuilder) WhereID(id string) *queryBuilder {
	b.id = id
	return b
}

func (b *queryBuilder) Find() (*Instance, error) {
	output, err := b.client.FindWorkflowInstance(b.workflowDefinition, b.id)
	if err != nil {
		return nil, err
	}

	properties := output["data"]["properties"]
	name := output["data"]["name"]

	return Manager.GetInstance(name, properties)
}

func (b *queryBuilder) Send(eventName string, eventData interface{}) {
	b.client.SendEvent(b.workflowDefinition, b.id, eventName, eventData)
}

// Kill a workflowDef instance
func (b *queryBuilder) Kill() (*queryBuilder, error) {
	err := b.client.KillWorkflow(b.workflowDefinition, b.id)
	return b, err
}

// Pause a workflowDef instance
func (b *queryBuilder) Pause() (*queryBuilder, error) {
	err := b.client.PauseWorkflow(b.workflowDefinition, b.id)
	return b, err
}

// Resume a workflowDef instance
func (b *queryBuilder) Resume() (*queryBuilder, error) {
	err := b.client.ResumeWorkflow(b.workflowDefinition, b.id)
	return b, err
}
