package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var TaskD = task.New(&D{})

type D struct{}

func (a *D) Handle() (interface{}, error) {
	time.Sleep(1 * time.Second)
	fmt.Println("Task D")
	return nil, nil
}
