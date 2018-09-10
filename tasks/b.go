package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var TaskB = task.New("TaskB", &B{})

type B struct{}

func (a *B) Handle() (interface{}, error) {
	time.Sleep(3 * time.Second)
	fmt.Println("Task B")
	return nil, nil
}
