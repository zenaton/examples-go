package workflow

/// The only requirements to write a workflow are:
//		1) Using either of the provided functions workflow.New or workflow.NewCustom;
//		2) Being idempotent. In more practical terms, it means it must implement a logical flow and NOT the tasks themselves.
//
// Idempotence implies that any actions (such as requesting a database, writing/reading a file, using current time,
// sending an email, echoing in console, etc.) that have side effects or that need access to potentially changing
// information MUST be done within tasks (not from within workflows).
//
// As Zenaton engine triggers the execution of the class describing a workflow each time it has to decide what to do
// next, failing to follow the idempotence requirement will lead to multiple executions of actions wrongly present in it.
//
// The provided method Dispatch is internally implemented to ensure idempotency.

import (
	"fmt"
	"reflect"

	"encoding/json"

	"github.com/zenaton/zenaton-go/v1/zenaton/internal/engine"
)

// Definition is the workflow definition. From a definition, you can create workflow instances with *Definition.New().
// Any two workflow Definitions cannot have the same name.
type Definition struct {
	name            string
	defaultInstance *Instance
	initFunc        reflect.Value
}

// New is the simpler way to create a new workflow Definition. You must provide a name and a handle function of the form:
// func () (interface{}, error). New is meant for simpler workflow definitions that don't need to keep track of state, have event handlers etc...
// For more options, use NewCustom instead.
// For example:
//
//			var SimpleWorkflow = workflow.New("SimpleWorkflow",
//				func() (interface{}, error) {
//					... // business logic of the workflow
//				})
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

// NewCustom creates a new workflow Definition. It takes a name and an instance of your type that implements the Handler Interface.
// The Handler interface has one method: Handle() (interface{}, error).
// For a simpler way to create a workflow Definition, use New.
//
// You can also optionally provide other methods:
// 		1) an Init method
// 			The Init method takes any number and type of arguments and initializes the workflow with data. This Init function will be
// 			called with the arguments passed to *Instance.New()
// 		2) an ID() string method
//			If you want to reference this workflow later, just implement a ID public method in your workflow implementation
// 			that provides the id Zenaton should use to retrieve this workflow instance, eg:
//					...
//					func (w *Welcome) ID() string {
//						return self.email
//					}
//					...
//
//			To be valid, this ID method MUST be unique (meaning in the same environment, you can not have two running
// 			instances of the same workflow with the same id). This ID method must have the form: func ID() string
// 		3) or an OnEvent() method
//			Workflows instances handle events an OnEvent method that will receive the event object as a parameter. Eg.
//				...
//				func (w *Welcome) OnEvent(name string, event interface{}){
//					if name == "AddressUpdatedEvent"{
//						eventMap = event.(map[string]interface{})
//						w.Address = eventMap["address"].(string)
//					}
//				}
//				...
//
//			The OnEvent method is called as soon the event is sent and a agent is available to execute it.
//			The workflow implementation MUST be idempotent. So the constraints on the OnEvent method are the same as the
//          handle method (it must implement a logical flow and NOT the tasks themselves.)
//			Note: an event is marshaled into and unmarshaled from json. This means that an event will contain the default
//          unmarshaled json types. The default unmarshaled type for structs or maps is map[string]interface{}. You can handle non-default types by sending the event as a json-encoded string and unmarshaling it in the OnEvent functio
//
// For example:
//
//		var WelcomeWorkflow = workflow.NewCustom("WelcomeWorkflow", &Welcome{})
//
//		type Welcome struct {
//				//Fields must be exported, as they will need to be serialized
//				Email string
//				SlackID string
//		}
//
//		func (w *Welcome) Init(user User) {
//			w.Email = user.Email
//			w.SlackID = user.SlackID
//		}
//
//		func (w *Welcome) Handle() (interface{}, error) {
//			SendWelcomeEmail(w.Email).Execute()
//			IntroduceUserThroughSlack(w.SlackID).Execute()
//		}
func NewCustom(name string, h engine.Handler) *Definition {

	rv := reflect.ValueOf(h)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic("must pass a pointer to NewCustom")
	}

	validateHandler(h)

	def := Definition{
		name:            name,
		defaultInstance: newInstance(name, h),
	}

	initFunc, ok := validateInit(h)
	if ok {
		def.initFunc = initFunc
	}

	UnsafeManager.setDefinition(def.name, &def)
	return &def
}

func validateHandler(value interface{}) {

	name := reflect.Indirect(reflect.ValueOf(value)).Type().Name()

	jsonV, err := json.Marshal(value)
	if err != nil {
		panic("workflow: Handler type '" + name + "' must be able to be marshaled to json. " + err.Error())
	}

	newV := reflect.New(reflect.TypeOf(value)).Interface()

	err = json.Unmarshal(jsonV, newV)
	if err != nil {
		panic("workflow: Handler type '" + name + "' must be able to be unmarshaled from json. " + err.Error())
	}
}

// WhereID takes an id (of a workflow instance) and returns a QueryBuilder.
// The QueryBuilder allows you to Find, Kill, Pause, and Resume workflow instances by id. You can also Send an event
// to a workflow.
func (d *Definition) WhereID(id string) *QueryBuilder {
	return newBuilder(d.name).whereID(id)
}

// Instance represents a workflow instance. You create new workflow instances by using *Definition.New(). Instances
// are runnable with Dispatch().
type Instance struct {
	name string
	engine.Handler
	OnEventer
	canonical string
	id        string
}

type OnEventer interface{OnEvent(string, interface{})}

// New returns an Instance. You must first have a workflow definition (created with New or NewCustom). If your Handler
// implementation has an Init() method, you can pass arguments to New which will then be passed to the Init() method.
func (d *Definition) New(args ...interface{}) *Instance {

	if len(args) > 0 {
		if !d.initFunc.IsValid() {
			panic("workflow: no Init() method set on: " + d.name)
		}

		values := []reflect.Value{reflect.ValueOf(d.defaultInstance.Handler)}
		for _, arg := range args {
			values = append(values, reflect.ValueOf(arg))
		}

		//this will panic if the arguments passed to New() don't match the provided Init function.

		d.callInit(args)

		jsondefaultHandler, err := json.Marshal(d.defaultInstance.Handler)
		if err != nil {
			panic(fmt.Sprint("workflow: must be able to json marshal the handler type... ", err.Error()))
		}

		newH := reflect.New(reflect.TypeOf(d.defaultInstance.Handler)).Interface()
		err = json.Unmarshal(jsondefaultHandler, &newH)

		if err != nil {
			panic(fmt.Sprint("workflow: must be able to json unmarshal into the handler type... ", err.Error()))
		}
	}
	return d.defaultInstance
}

func (d *Definition) callInit(args []interface{}) {
	//here we recover the panic just to add some more helpful information, then we re-panic
	defer func() {
		r := recover()
		if r != nil {
			panic(fmt.Sprint("workflow: arguments passed to Definition.New() must be of the same type and quantity of those defined in the Init function... ", r))
		}
	}()

	values := []reflect.Value{reflect.ValueOf(d.defaultInstance.Handler)}
	for _, arg := range args {
		values = append(values, reflect.ValueOf(arg))
	}

	//this will panic if the arguments passed to New() don't match the provided Init function.
	d.initFunc.Call(values)
}

// Dispatch launches a workflow asynchronously.
func (i *Instance) Dispatch() {
	e := engine.NewEngine()
	e.Dispatch([]engine.Job{i})
}

func validateInit(value interface{}) (reflect.Value, bool) {

	rt := reflect.TypeOf(value)

	initMethod, ok := rt.MethodByName("Init")
	if !ok {
		return reflect.Value{}, false
	}

	return initMethod.Func, true
}

func newInstance(name string, h engine.Handler) *Instance {

	i := &Instance{
		name:    name,
		Handler: h,
	}

	eventer, ok := h.(interface {
		OnEvent(string, interface{})
	})

	if ok {
		i.OnEventer = eventer
	}

	return i
}

// GetData allows you to retrieve the underlying Handler implementation from a workflow Instance
// For example:
// 			WelcomeFlow.WhereID(email).Find().GetData().(*Welcome)
// GetData returns a Handler interface, so you can do a type assertion to your base Type if you like.
// then you can inspect the data associated with your Welcome type (in this example).
func (i Instance) GetData() engine.Handler { return i.Handler }

// GetName retrieves the name of an Instance. This is used in the agent code, so must be exported.
func (i Instance) GetName() string { return i.name }

// LaunchInfo is needed for the agent. You shouldn't need to use this.
func (i Instance) LaunchInfo() engine.LaunchInfo {
	return engine.LaunchInfo{
		Type:      "workflow",
		Name:      i.name,
		Canonical: i.canonical,
		ID:        i.GetCustomID(),
		Data:      i.Handler,
	}
}

// GetCustomID retrieves an Instance ID. This will be "" if you don't have a ID() string method in your workflow
func (i *Instance) GetCustomID() string {
	ider, ok := i.Handler.(interface{ID()string})
	if ok {
		return ider.ID()
	}
	return ""
}

// GetCanonical returns the versioned workflow name instead of the actual workflow name.
func (i *Instance) GetCanonical() string {
	return i.canonical
}
