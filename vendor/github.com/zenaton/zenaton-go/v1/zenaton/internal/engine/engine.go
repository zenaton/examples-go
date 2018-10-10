package engine

import (
	"github.com/zenaton/zenaton-go/v1/zenaton/internal/client"
)

var instance = &Engine{
	client: client.NewClient(false),
}

type Engine struct {
	client    *client.Client
	processor Processor
}

func NewEngine() *Engine {
	return instance
}

type Processor interface {
	Process([]Job, bool) ([]interface{}, []string, []error)
}

type LaunchInfo struct {
	Type      string
	Name      string
	Canonical string
	ID        string
	Data      interface{}
}

type Handler interface {
	Handle() (interface{}, error)
}

type Job interface {
	Handle() (interface{}, error)
	LaunchInfo() LaunchInfo
	GetName() string
	GetData() Handler
}

func (e *Engine) Execute(jobs []Job) ([]interface{}, []string, []error) {

	// local execution
	if e.processor == nil || len(jobs) == 0 {
		var outputs []interface{}
		var errs []error
		for _, job := range jobs {
			out, err := job.GetData().Handle()

			errs = append(errs, err)
			outputs = append(outputs, out)
		}

		return outputs, nil, errs
	}

	outputValues, serializedOutputs, errs := e.processor.Process(jobs, true)
	return outputValues, serializedOutputs, errs
}

func (e *Engine) Dispatch(jobs []Job) {

	if e.processor == nil || len(jobs) == 0 {

		for _, job := range jobs {
			li := job.LaunchInfo()
			if li.Type == "workflow" {
				client.NewClient(false).StartWorkflow(li.Name, li.Canonical, li.ID, li.Data)
			} else {
				job.Handle()
			}
		}

		return
	}

	e.processor.Process(jobs, false)
}

func (e *Engine) SetProcessor(processor Processor) {
	e.processor = processor
}
