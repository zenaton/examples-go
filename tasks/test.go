package tasks

import (
	"fmt"

	"github.com/zenaton/examples-go/tasks/log"
	"github.com/zenaton/zenaton-go/v1/zenaton/errors"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var TestTask = task.New("TestTask", &Test{})

type Test struct {
	Return   interface{}
	Error    string
	Panic    string
	Print    string
	Relaunch bool
}

func (t *Test) Handle() (interface{}, error) {
	if t.Print != "" {
		fmt.Println(t.Print)
	}

	if t.Panic != "" {
		panic(t.Panic)
	}

	if t.Error != "" {
		return nil, errors.New("testErrorName", t.Error)
	}

	if t.Relaunch {
		var out interface{}
		err := TestTask.NewInstance(&Test{
			Return: t.Return,
			Error:  t.Error,
		}).Execute(&out)
		return out, err
	}

	return t.Return, nil
}

var TaskRunnerTask = task.New("TaskRunnerTask", &TaskRunner{})

type TaskRunner struct{}

func (t *TaskRunner) Handle() (interface{}, error) {
	//without return
	err := TestTask.NewInstance().Execute()
	log.Println("err: ", err)
	if err != nil {
		return nil, err
	}

	//with return
	var out interface{}
	err = TestTask.NewInstance(&Test{
		Return: "test return",
	}).Execute(&out)

	return out, err
}
