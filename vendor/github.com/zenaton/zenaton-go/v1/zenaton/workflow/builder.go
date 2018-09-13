package workflow

import (
	"github.com/zenaton/zenaton-go/v1/zenaton/client"
)

type QueryBuilder struct {
	workflowDefinition string
	id                 string
	client             *client.Client
}

func NewBuilder(name string) *QueryBuilder {
	return &QueryBuilder{
		client:             client.NewClient(false),
		workflowDefinition: name,
	}
}

func (b *QueryBuilder) WhereID(id string) *QueryBuilder {
	b.id = id
	return b
}

func (b *QueryBuilder) Find() (*Instance, error) {
	output, err := b.client.FindWorkflowInstance(b.workflowDefinition, b.id)
	if err != nil {
		return nil, err
	}

	properties := output["data"]["properties"]
	name := output["data"]["name"]

	return Manager.GetInstance(name, properties)
}

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
