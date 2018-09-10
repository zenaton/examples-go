package workflows

import (
	"fmt"
	"github.com/zenaton/examples-go/tasks"
	"github.com/zenaton/examples-go/tasks/log"
	"github.com/zenaton/zenaton-go/v1/zenaton/errors"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var WithTaskWorkflow = workflow.New("WithTaskWorkflow", &WithTask{})

type WithTask struct {
	Task *task.Task
}

func (wt *WithTask) Handle() (interface{}, error) { return nil, nil }

var TestW = workflow.New("TestWorkflow", &Test2{})

type Test2 struct {
	Relaunch   bool
	Parallel   bool
	Print      string
	Return     string
	Error      string
	Panic      string
	TaskReturn interface{}
	TaskError  string
	TaskPanic  string
	TaskPrint  string
	IDstr      string
}

func (t Test2) Handle() (interface{}, error) {
	if t.Parallel {

		var t1Return tasks.Test
		var t2Return tasks.Test

		t1 := tasks.TestTask.NewInstance(&tasks.Test{
			Return: t.TaskReturn,
			Error:  t.TaskError,
			Panic:  t.TaskPanic,
			Print:  t.TaskPrint,
		})

		t2 := t1

		err := task.Parallel{t1, t2}.Execute(&t1Return, &t2Return)
		log.Println("out1: ", t1Return)
		log.Println("out2: ", t2Return)
		log.Println("err: ", err)

	} else {
		var out interface{}
		err := tasks.TestTask.NewInstance(&tasks.Test{
			Return:   t.TaskReturn,
			Error:    t.TaskError,
			Panic:    t.TaskPanic,
			Print:    t.TaskPrint,
			Relaunch: t.Relaunch,
		}).Execute(&out)

		log.Println("out: ", out)
		log.Println("err: ", err)
	}

	if t.Panic != "" {
		panic(t.Panic)
	}

	if t.Error != "" {
		return nil, errors.New("testWorkflowError", "testErrorMessage")
	}

	return t.Return, nil
}

func (t Test2) ID() string {
	if t.IDstr != "" {
		return t.IDstr
	}
	return "MyID"
}

var TestRelaunchTaskWorkflow = workflow.New("TestRelaunchTaskWorkflow", &TestRelaunchTask{})

type TestRelaunchTask struct{}

func (t TestRelaunchTask) Handle() (interface{}, error) {

	var out interface{}
	err := tasks.TaskRunnerTask.NewInstance().Execute(&out)

	log.Println("out: ", out)
	log.Println("err: ", err)

	return nil, nil
}

var TestEventValueWorkflow = workflow.New("TestEventValueWorkflow", &TestEventValue{})

type TestEventValue struct{}

func (tev TestEventValue) Handle() (interface{}, error) {

	event, err := task.Wait().ForEvent("MyEvent").Seconds(15).Execute()
	fmt.Println("wait for event: ", event, err)

	return nil, nil
}

func (tev *TestEventValue) OnEvent(eventName string, eventData interface{}) {
	fmt.Println("onEvent: ", eventName, eventData)
}

func (tev TestEventValue) ID() string {
	fmt.Println("bob")
	return "TestEventValueID"
}
