package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var TaskC = task.New(&C{})

type C struct{}

func (a *C) Handle() (interface{}, error) {
	time.Sleep(2 * time.Second)
	fmt.Println("Task C")
	return nil, nil
}
