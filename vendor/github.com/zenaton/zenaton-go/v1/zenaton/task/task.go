package task

import (
	"reflect"

	"encoding/json"

	"fmt"

	"errors"
	"github.com/zenaton/zenaton-go/v1/zenaton/engine"
	"github.com/zenaton/zenaton-go/v1/zenaton/interfaces"
	"github.com/zenaton/zenaton-go/v1/zenaton/service/serializer"
)

type Task struct {
	name string
	interfaces.Handler
}

type taskType struct {
	name        string
	defaultTask *Task
}

type defaultHandler struct {
	handlerFunc func() (interface{}, error)
}

func (dh *defaultHandler) Handle() (interface{}, error) {
	return dh.handlerFunc()
}

func NewDefault(name string, handlerFunc func() (interface{}, error)) *taskType {
	return New(name, &defaultHandler{
		handlerFunc: handlerFunc,
	})
}

func New(name string, h interfaces.Handler) *taskType {
	rv := reflect.ValueOf(h)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic("must pass a pointer to NewWorkflow")
	}

	validateHandler(h)

	taskT := taskType{
		name:        name,
		defaultTask: newInstance(name, h),
	}

	NewTaskManager().setClass(taskT.name, &taskT)
	return &taskT
}

func (tt *taskType) NewInstance(handlers ...interfaces.Handler) *Task {
	if len(handlers) > 1 {
		panic("must only pass one handler to taskType.NewInstance()")
	}

	if len(handlers) == 1 {
		h := handlers[0]
		validateHandler(h)
		return newInstance(tt.name, h)
	} else {
		return tt.defaultTask
	}
}

func validateHandler(value interface{}) {

	name := reflect.Indirect(reflect.ValueOf(value)).Type().Name()

	jsonV, err := json.Marshal(value)
	if err != nil {
		panic("handler type '" + name + "' must be able to be marshaled to json. " + err.Error())
	}

	newV := reflect.New(reflect.TypeOf(value)).Interface()

	err = json.Unmarshal(jsonV, newV)
	if err != nil {
		panic("handler type '" + name + "' must be able to be unmarshaled from json. " + err.Error())
	}
}

func newInstance(name string, h interfaces.Handler) *Task {
	return &Task{
		name:    name,
		Handler: h,
	}
}

func (t Task) GetName() string { return t.name }

func (t Task) GetData() interface{} { return t.Handler }

func (t Task) Async() error {
	t.Handler.Handle()
	return nil
}

type MaxProcessingTimer interface {
	MaxTime() int64
}

func (t Task) MaxProcessingTime() int64 {
	maxer, ok := t.Handler.(MaxProcessingTimer)
	if ok {
		return maxer.MaxTime()
	}
	return -1
}

type taskExecution struct {
	outputValue     interface{}
	serializedValue string
	err             error
}

func (te *taskExecution) Output(values ...interface{}) error {

	if len(values) > 1 {
		panic("must pass a maximum of 1 value to Output")
	}

	if te.serializedValue != "" {
		var value interface{}
		if len(values) == 1 {
			value = values[0]
		}
		return outputFromSerialized(value, te.serializedValue)

	} else {

		if len(values) == 1 {
			value := values[0]
			outputFromInterface(value, te.outputValue)
		}

		return te.err
	}
}

func outputFromInterface(to, from interface{}) {
	rv := reflect.ValueOf(to)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic(fmt.Sprint("must pass a non-nil pointer to task.Execute"))
	}

	if from != nil && to != nil {
		outV := reflect.ValueOf(from)
		if outV.IsValid() {
			rv.Elem().Set(outV)
		}
	}
}

func outputFromSerialized(to interface{}, from string) error {
	var combinedOutput map[string]json.RawMessage

	err := serializer.Decode(from, &combinedOutput)
	if err != nil {
		panic(err)
	}

	if to != nil {
		err = serializer.Decode(string(combinedOutput["output"]), to)
		if err != nil {
			panic(err)
		}
	}

	if combinedOutput["error"] != nil {
		return errors.New(string(combinedOutput["error"]))
	}
	return nil
}

func (t *Task) Execute() *taskExecution {

	outputValues, serializedValues, errs := engine.NewEngine().Execute([]interfaces.Job{t})

	var ex taskExecution

	if outputValues != nil {
		ex.outputValue = outputValues[0]
		ex.err = errs[0]
	}

	if serializedValues != nil {
		ex.serializedValue = serializedValues[0]
	}

	return &ex
}

func (t *Task) Dispatch() {
	e := engine.NewEngine()
	e.Dispatch([]interfaces.Job{t})
}

type parallelExecution struct {
	outputValues     []interface{}
	serializedValues []string
	errors           []error
}

func (pe *parallelExecution) Output(values ...interface{}) []error {

	if len(values) != len(pe.outputValues) && len(values) != len(pe.serializedValues) {
		panic(fmt.Sprint("task: number of parallel tasks and return value pointers do not match"))
	}

	if len(values) == 0 {
		values = make([]interface{}, len(pe.errors))
	}

	var errs []error

	if pe.serializedValues != nil {
		for i := range pe.serializedValues {
			err := outputFromSerialized(values[i], pe.serializedValues[i])
			errs = append(errs, err)
		}
	} else {

		for i := range pe.outputValues {
			if values[i] != nil {
				value := values[0]
				outputFromInterface(value, pe.outputValues[i])
			}
		}

		errs = pe.errors
	}

	for _, e := range errs {
		if e != nil {
			return pe.errors
		}
	}
	return nil
}

type Parallel []*Task

func (ts Parallel) Dispatch() {
	e := engine.NewEngine()
	var jobs []interfaces.Job
	for _, task := range ts {
		jobs = append(jobs, task)
	}
	e.Dispatch(jobs)
}

func (ts Parallel) Execute() *parallelExecution {

	////todo: test this
	//if len(returnValues) > 0 {
	//	if len(returnValues) != len(ts) {
	//		panic(fmt.Sprint("task: number of parallel tasks (", strconv.Itoa(len(ts)),
	//			") and return value pointers (", strconv.Itoa(len(returnValues)), ") do not match"))
	//	}
	//	for i, returnValue := range returnValues {
	//		rv := reflect.ValueOf(returnValue)
	//		if rv.Kind() != reflect.Ptr || rv.IsNil() {
	//			panic(fmt.Sprint("item at index ", i, " must pass a non-nil pointer to task.Execute"))
	//		}
	//	}
	//} else {
	//	returnValues = make([]interface{}, len(ts))
	//}

	e := engine.NewEngine()
	var jobs []interfaces.Job
	for _, task := range ts {
		jobs = append(jobs, task)
	}
	values, serializedValues, errors := e.Execute(jobs)

	return &parallelExecution{
		outputValues:     values,
		serializedValues: serializedValues,
		errors:           errors,
	}
}
