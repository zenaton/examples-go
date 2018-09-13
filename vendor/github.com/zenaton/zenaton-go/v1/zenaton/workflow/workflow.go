package workflow

import (
	"reflect"

	"encoding/json"

	"github.com/zenaton/zenaton-go/v1/zenaton/client"
	"github.com/zenaton/zenaton-go/v1/zenaton/engine"
	"github.com/zenaton/zenaton-go/v1/zenaton/interfaces"
)

type Definition struct {
	name            string
	defaultInstance *Instance
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

	Manager.setDefinition(def.name, &def)
	return &def
}

func (d *Definition) WhereID(id string) *QueryBuilder {
	return NewBuilder(d.name).WhereID(id)
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

func (d *Definition) NewInstance(handlers ...interfaces.Handler) *Instance {

	if len(handlers) > 1 {
		panic("must only pass one handler to Definition.NewInstance()")
	}

	if len(handlers) == 1 {
		h := handlers[0]
		validateHandler(h)
		return newInstance(d.name, h)
	} else {
		return d.defaultInstance
	}
}

func (i *Instance) Dispatch() error {
	e := engine.NewEngine()
	return e.Dispatch([]interfaces.Job{i})
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
