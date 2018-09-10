package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var TaskA = task.New("TaskA", &A{})

type A struct {
	Field string
}

func (a *A) Handle() (interface{}, error) {
	time.Sleep(4 * time.Second)
	fmt.Println("Task A")
	return nil, nil
}
