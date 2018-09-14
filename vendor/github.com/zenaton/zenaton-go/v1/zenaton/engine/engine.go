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
	Process([]interfaces.Job, bool) ([]interface{}, []string, error)
}

func (e *Engine) Execute(jobs []interfaces.Job) ([]interface{}, []string) {

	// local execution
	if e.processor == nil || len(jobs) == 0 {
		var outputs []interface{}
		for _, job := range jobs {
			out, err := job.Handle()
			if err != nil {
				panic(err)
			}
			outputs = append(outputs, out)
		}

		return outputs, nil
	}

	outputValues, serializedOutputs, err := e.processor.Process(jobs, true)
	if err != nil {
		panic(err)
	}

	return outputValues, serializedOutputs
}

func (e *Engine) Dispatch(jobs []interfaces.Job) {

	if e.processor == nil || len(jobs) == 0 {

		for _, job := range jobs {
			err := job.Async()
			if err != nil {
				panic(err)
			}
		}

		return
	}

	_, _, err := e.processor.Process(jobs, false)
	if err != nil {
		panic(err)
	}
}

func (e *Engine) SetProcessor(processor Processor) {
	e.processor = processor
}
