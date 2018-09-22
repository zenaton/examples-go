package workflow

import (
	fmt "fmt"
	"reflect"

	"encoding/json"

	"github.com/zenaton/zenaton-go/v1/zenaton/client"
	"github.com/zenaton/zenaton-go/v1/zenaton/engine"
	"github.com/zenaton/zenaton-go/v1/zenaton/interfaces"
)

type Definition struct {
	name            string
	defaultInstance *Instance
	initFunc        reflect.Value
}

func New(name string, handlerFunc func() (interface{}, error)) *Definition {
	return NewCustom(name, &defaultHandler{
		handlerFunc: handlerFunc,
	})
}

func NewCustom(name string, h interfaces.Handler) *Definition {

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

	Manager.setDefinition(def.name, &def)
	return &def
}

func (d *Definition) WhereID(id string) *queryBuilder {
	return newBuilder(d.name).WhereID(id)
}

type defaultHandler struct {
	handlerFunc func() (interface{}, error)
}

func (dh *defaultHandler) Handle() (interface{}, error) {
	return dh.handlerFunc()
}

type Instance struct {
	name string
	interfaces.Handler
	interfaces.OnEventer
	canonical string
	id        string
}

func (d *Definition) New(args ...interface{}) *Instance {

	if len(args) > 0 {
		if !d.initFunc.IsValid() {
			panic("workflow: no Init() method set on: " + d.name)
		}

		// here we recover the panic just to add some more helpful information, then we re-panic
		defer func() {
			r := recover()
			if r != nil {
				panic(fmt.Sprint("workflow: arguments passed to Definition.New() must be of the same time and quantity of those defined in the Init function"))
			}
		}()

		values := []reflect.Value{reflect.ValueOf(d.defaultInstance.Handler)}
		for _, arg := range args {
			values = append(values, reflect.ValueOf(arg))
		}

		//this will panic if the arguments passed to New() don't match the provided Init function.
		d.initFunc.Call(values)
	}
	return d.defaultInstance
}

func validateInit(value interface{}) (reflect.Value, bool) {

	rt := reflect.TypeOf(value)

	initMethod, ok := rt.MethodByName("Init")
	if !ok {
		return reflect.Value{}, false
	}

	return initMethod.Func, true
}

func (i *Instance) Dispatch() {
	e := engine.NewEngine()
	e.Dispatch([]interfaces.Job{i})
}

func newInstance(name string, h interfaces.Handler) *Instance {

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

func (i Instance) GetName() string { return i.name }

func (i Instance) GetData() interface{} { return i.Handler }

func (i Instance) Async() error {
	return client.NewClient(false).StartWorkflow(i.name, i.canonical, i.GetCustomID(), i.Handler)
}

func (i *Instance) GetCustomID() string {
	ider, ok := i.Handler.(interfaces.IDer)
	if ok {
		return ider.ID()
	}
	return ""
}

func (i *Instance) GetCanonical() string {
	return i.canonical
}

func (i *Instance) setCanonical(canonical string) {
	i.canonical = canonical
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
