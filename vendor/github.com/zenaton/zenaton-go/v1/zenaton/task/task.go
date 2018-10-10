package task

import (
	"reflect"

	"encoding/json"

	"fmt"

	"errors"

	"github.com/zenaton/zenaton-go/v1/zenaton/internal/engine"
	"github.com/zenaton/zenaton-go/v1/zenaton/service/serializer"
)

// Definition is the task definition. From a definition, you can create task instances with *Definition.New().
// Any two task Definitions cannot have the same name.
type Definition struct {
	name        string
	defaultTask *Instance
	initFunc    reflect.Value
}

// New is the simpler way to create a new task Definition. You must provide a name and a handle function of the form:
// func () (interface{}, error). New is meant for simpler task definitions that don't need an Init() function,
// or a custom Handler implementation
// For more options, use NewCustom instead.
//
// For example:
//
//		var SimpleTask = task.New("SimpleTask",
//    		func() (interface{}, error) {
//        		... // business logic of the task
//    		})
func New(name string, handlerFunc func() (interface{}, error)) *Definition {
	return NewCustom(name, &defaultHandler{
		handlerFunc: handlerFunc,
	})
}

type defaultHandler struct {
	handlerFunc func() (interface{}, error)
}

func (dh *defaultHandler) Handle() (interface{}, error) {
	return dh.handlerFunc()
}

// NewCustom creates a new task Definition. It takes a name and an instance of your type that implements the Handler Interface.
// The Handler interface has one method: Handle() (interface{}, error).
//
// You can optionally provide :
// 		1) an Init Method.
// 			The Init method can take any number and type of arguments and initialize the
// 			task with data. This Init function will be called with the arguments passed to *Instance.New()
//		2) MaxTime() int64
//        	When a task fails, an exception is usually thrown and the Zenaton engine is alerted. But in some rare cases
// 			(such as an agent crash), it is not possible to know for sure that something went wrong. To handle this
// 			situation, a task will be considered a failure if its execution lasts more than 30 seconds. You can modify
// 			this value by using a MaxTime method
//
// 			Your MaxTime() method must have exactly this signature. The returned int64 will be interpreted as the number
// 			of seconds before a task is considered timed out.
//
//			MaxTime() will not be used to actually stop a task from running after the given time. Instead, when the time
//			is reached, your zenaton interface (https://zenaton.com/app/monitoring) will show a timeout error for this
// 			task, and you can retry/kill the task if you wish.
//
// For a simpler way to create a task Definition, use New.
//
// For example:
//
// 		var CustomTypeTask = task.NewCustom("CustomTypeTask", &CustomType{})
//
//		type CustomType struct {
//			// Fields must be exported, as they will need to be serialized
//			Price int
//			Key string
//		}
//
//		func (ct *CustomType) Init(price int, key string) {
//			ct.Price = 2 * price
//			ct.Key = "key" + key
//		}
//
//		func (ct *CustomType) Handle() (interface{}, error) {
//			... //  task implementation that can now use ct.Price and ct.Key
//		}
//
//		func (ct *CustomType) MaxTime() int64 {
//			return 180 //after 3 minutes this task will be considered a failure
//		}
func NewCustom(name string, h engine.Handler) *Definition {
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

	UnsafeManager.setDefinition(taskT.name, &taskT)
	return &taskT
}

func newInstance(name string, h engine.Handler) *Instance {
	return &Instance{
		name:    name,
		Handler: h,
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

func validateInit(value interface{}) (reflect.Value, bool) {

	rt := reflect.TypeOf(value)

	initMethod, ok := rt.MethodByName("Init")
	if !ok {
		return reflect.Value{}, false
	}

	return initMethod.Func, true
}

// Instance represents a task instance. You create new task instances by using *Definition.New(). Instances
// are runnable with Dispatch().
type Instance struct {
	name string
	engine.Handler
}

// New returns an Instance. You must first have a task definition (created with New or NewCustom). If your Handler
// implementation has an Init() method, you can pass arguments to New which will then be passed to the Init() method.
func (tt *Definition) New(args ...interface{}) *Instance {

	if len(args) > 0 {
		if !tt.initFunc.IsValid() {
			panic("task: no Init() method set on: " + tt.name)
		}

		tt.callInit(args)
	}

	jsonDefaultHandler, err := json.Marshal(tt.defaultTask.Handler)
	if err != nil {
		panic("task: must be able to marshal handler to json: " + err.Error())
	}

	newH := reflect.New(reflect.TypeOf(tt.defaultTask.Handler)).Interface()
	err = json.Unmarshal(jsonDefaultHandler, &newH)
	if err != nil {
		panic(fmt.Sprint("task: must be able to json unmarshal into the handler type... ", err.Error()))
	}

	return tt.defaultTask
}

func (tt *Definition) callInit(args []interface{}) {
	//here we recover the panic just to add some more helpful information, then we re-panic
	defer func() {
		r := recover()
		if r != nil {
			panic(fmt.Sprint("task: arguments passed to Definition.New() must be of the same type and quantity of those defined in the Init function... ", r))
		}
	}()

	values := []reflect.Value{reflect.ValueOf(tt.defaultTask.Handler)}
	for _, arg := range args {
		values = append(values, reflect.ValueOf(arg))
	}

	//this will panic if the arguments passed to New() don't match the provided Init function.
	tt.initFunc.Call(values)
}

// Dispatch launches a task instance asynchronously.
func (i *Instance) Dispatch() {
	e := engine.NewEngine()
	e.Dispatch([]engine.Job{i})
}

// Execute launches a task instance synchronously (the workflow will block until this task is done).
// Execute returns a Execution, which you can use to get the output and error of the task.
// for example:
//
// 		var a int
//		err := tasks.A.New().Execute().Output(&a)
// 		if err != nil {
//    		... //handle error
//		}
//
// Note: If you have a custom error type, the information will be lost. Here we just return a standard go error
// where err.Error() matches the output of the err.Error() that was returned from the task.
func (i *Instance) Execute() Execution {

	outputValues, serializedValues, errs := engine.NewEngine().Execute([]engine.Job{i})

	var ex Execution

	if outputValues != nil {
		ex.outputValue = outputValues[0]
		ex.err = errs[0]
	}

	if serializedValues != nil {
		ex.serializedValue = serializedValues[0]
	}

	return ex
}

// Execution represents the output and error of the task.
type Execution struct {
	outputValue     interface{}
	serializedValue string
	err             error
}

// Output allows you to to get the output and error of the task.
// for example:
//
// 		var a int
//		err := tasks.A.New().Execute().Output(&a)
// 		if err != nil {
//    		... //handle error
//		}
//
// Note: If you have a custom error type, the information will be lost. Here we just return a standard go error
// where err.Error() matches the output of the err.Error() that was returned from the task.
func (te Execution) Output(values ...interface{}) error {

	if len(values) > 1 {
		panic("must pass a maximum of 1 value to Output")
	}

	if te.serializedValue != "" {
		var value interface{}
		if len(values) == 1 {
			value = values[0]
		}
		return outputFromSerialized(value, te.serializedValue)

	}

	if len(values) == 1 {
		value := values[0]
		outputFromInterface(value, te.outputValue)
	}

	return te.err
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

// Parallel is just a slice of *Instances that can be run in parallel with Execute() or Dispatch().
type Parallel []*Instance

// Execute will execute a Parallel (a slice of Instances) in parallel and wait for their completion.
// Execute returns a ParallelExecution which can be used to get the outputs and errors of the tasks executed.
//
// For parallel tasks, you will receive a slice of errors. This slice will be nil if no error occurred. If there was
// an error in one of the parallel tasks, you will receive a slice of the same length as the input tasks, and the index
// of the task that produced an error will be the same index as the non-nil err in the slice of errors
//
//  var a int
//	var b int
//
//	errs := task.Parallel{
//	    tasks.A.New(),
//	    tasks.B.New(),
//	}.Execute().Output(&a, &b)
//
//	if errs != nil {
//	    if errs[0] != nil {
//	        //  tasks.A error
//	    }
//	    if errs[1] != nil {
//	        //  tasks.B error
//	    }
//	}
//
// Here, tasks A and B will be executed in parallel, and we wait for all of them to end before continuing. You can
// retrieve the outputs of these tasks by passing pointers to .Output()
func (ts Parallel) Execute() ParallelExecution {

	e := engine.NewEngine()
	var jobs []engine.Job
	for _, task := range ts {
		jobs = append(jobs, task)
	}
	values, serializedValues, errors := e.Execute(jobs)

	return ParallelExecution{
		outputValues:     values,
		serializedValues: serializedValues,
		errors:           errors,
	}
}

// Dispatch will launch the the tasks in parallel and not wait for them to complete before moving on. Thus:
//
//		task.Parallel{
//			tasks.A.New(),
//			tasks.B.New(),
//		}.Dispatch()
//
// should be equivalent to:
//
// 		tasks.A.New().Dispatch()
// 		tasks.B.New().Dispatch()
//
func (ts Parallel) Dispatch() {
	e := engine.NewEngine()
	var jobs []engine.Job
	for _, task := range ts {
		jobs = append(jobs, task)
	}
	e.Dispatch(jobs)
}

// ParallelExecution represents the outputs and errors of the Parallel tasks.
// To get the output, use ParallelExecution.Output()
type ParallelExecution struct {
	outputValues     []interface{}
	serializedValues []string
	errors           []error
}

// Output gets the output of a Parallel execution
//
//
// For parallel tasks, you will receive a slice of errors. This slice will be nil if no error occurred. If there was
// an error in one of the parallel tasks, you will receive a slice of the same length as the input tasks, and the index
// of the task that produced an error will be the same index as the non-nil err in the slice of errors
//
//  var a int
//	var b int
//
//	errs := task.Parallel{
//	    tasks.A.New(),
//	    tasks.B.New(),
//	}.Execute().Output(&a, &b)
//
//	if errs != nil {
//	    if errs[0] != nil {
//	        //  tasks.A error
//	    }
//	    if errs[1] != nil {
//	        //  tasks.B error
//	    }
//	}
//
// Here, tasks A and B will be executed in parallel, and we wait for all of them to end before continuing. You can
// retrieve the outputs of these tasks by passing pointers to .Output()
func (pe ParallelExecution) Output(values ...interface{}) []error {

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

// GetName simply returns the name of the Instance
func (i *Instance) GetName() string { return i.name }

// GetData allows you to retrieve the underlying handler implementation of a task Instance
func (i *Instance) GetData() engine.Handler { return i.Handler }

// LaunchInfo returns some information about what type of Instance you have (either a task or a workflow).
func (i *Instance) LaunchInfo() engine.LaunchInfo {
	return engine.LaunchInfo{
		Type: "task",
	}
}
