package engine

import (
	"github.com/zenaton/zenaton-go/v1/zenaton/client"
	"github.com/zenaton/zenaton-go/v1/zenaton/interfaces"
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
	Process([]interfaces.Job, bool) ([]interface{}, []string, []error)
}

func (e *Engine) Execute(jobs []interfaces.Job) ([]interface{}, []string, []error) {

	// local execution
	if e.processor == nil || len(jobs) == 0 {
		var outputs []interface{}
		var errs []error
		for _, job := range jobs {
			out, err := job.Handle()

			errs = append(errs, err)
			outputs = append(outputs, out)
		}

		return outputs, nil, errs
	}

	outputValues, serializedOutputs, errs := e.processor.Process(jobs, true)
	return outputValues, serializedOutputs, errs
}

func (e *Engine) Dispatch(jobs []interfaces.Job) []error {

	if e.processor == nil || len(jobs) == 0 {

		var errs []error
		for _, job := range jobs {
			err := job.Async()
			errs = append(errs, err)
		}

		return errs
	}

	_, _, errs := e.processor.Process(jobs, false)
	return errs
}

func (e *Engine) SetProcessor(processor Processor) {
	e.processor = processor
}
