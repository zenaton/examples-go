package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var TaskC = task.New2(&C{})

type C struct{}

func (a *C) Handle() {
	time.Sleep(2 * time.Second)
	fmt.Println("Task C")
}
