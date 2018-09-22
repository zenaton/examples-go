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

type Instance struct {
	name string
	interfaces.Handler
}

type Definition struct {
	name        string
	defaultTask *Instance
	initFunc    reflect.Value
}

type defaultHandler struct {
	handlerFunc func() (interface{}, error)
}

func (dh *defaultHandler) Handle() (interface{}, error) {
	return dh.handlerFunc()
}

func New(name string, handlerFunc func() (interface{}, error)) *Definition {
	return NewCustom(name, &defaultHandler{
		handlerFunc: handlerFunc,
	})
}

func NewCustom(name string, h interfaces.Handler) *Definition {
	rv := reflect.ValueOf(h)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic("must pass a pointer to NewWorkflow")
	}

	validateHandler(h)

	taskT := Definition{
		name:        name,
		defaultTask: newInstance(name, h),
	}
	initFunc, ok := validateInit(h)
	if ok {
		taskT.initFunc = initFunc
	}

	Manager.setDefinition(taskT.name, &taskT)
	return &taskT
}

func (tt *Definition) New(args ...interface{}) *Instance {

	if len(args) > 0 {
		if !tt.initFunc.IsValid() {
			panic("task: no Init() method set on: " + tt.name)
		}

		//here we recover the panic just to add some more helpful information, then we re-panic
		defer func() {
			r := recover()
			if r != nil {
				panic(fmt.Sprint("task: arguments passed to Definition.New() must be of the same time and quantity of those defined in the Init function"))
			}
		}()

		values := []reflect.Value{reflect.ValueOf(tt.defaultTask.Handler)}
		for _, arg := range args {
			values = append(values, reflect.ValueOf(arg))
		}

		//this will panic if the arguments passed to New() don't match the provided Init function.
		tt.initFunc.Call(values)
	}
	return tt.defaultTask
}

func validateInit(value interface{}) (reflect.Value, bool) {

	rt := reflect.TypeOf(value)

	initMethod, ok := rt.MethodByName("Init")
	if !ok {
		return reflect.Value{}, false
	}

	return initMethod.Func, true
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

func newInstance(name string, h interfaces.Handler) *Instance {
	return &Instance{
		name:    name,
		Handler: h,
	}
}

func (i *Instance) GetName() string { return i.name }

func (i *Instance) GetData() interface{} { return i.Handler }

func (i *Instance) Async() error {
	i.Handler.Handle()
	return nil
}

type MaxProcessingTimer interface {
	MaxTime() int64
}

func (i *Instance) MaxProcessingTime() int64 {
	maxer, ok := i.Handler.(MaxProcessingTimer)
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
		panic(fmt.Sprint("must pass a non-nil pointer to task.Output"))
	}

	if from != nil && to != nil {
		outV := reflect.ValueOf(from)
		if outV.IsValid() {
			rv.Elem().Set(outV)
		}
	}
}

func outputFromSerialized(to interface{}, from string) error {

	rv := reflect.ValueOf(to)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic(fmt.Sprint("must pass a non-nil pointer to task.Output"))
	}

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

func (i *Instance) Execute() *taskExecution {

	outputValues, serializedValues, errs := engine.NewEngine().Execute([]interfaces.Job{i})

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

func (i *Instance) Dispatch() {
	e := engine.NewEngine()
	e.Dispatch([]interfaces.Job{i})
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

type Parallel []*Instance

func (ts Parallel) Dispatch() {
	e := engine.NewEngine()
	var jobs []interfaces.Job
	for _, task := range ts {
		jobs = append(jobs, task)
	}
	e.Dispatch(jobs)
}

func (ts Parallel) Execute() *parallelExecution {

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
