package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var TaskB = task.New2(&B{})

type B struct{}

func (a *B) Handle() {
	time.Sleep(3 * time.Second)
	fmt.Println("Task B")
}
