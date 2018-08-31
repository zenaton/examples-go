package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

func init() {
	//todo: call this registerTask
	task.New(&C{})
}

type C struct{}

func (a *C) Handle() {
	time.Sleep(2 * time.Second)
	fmt.Println("Task C")
}
