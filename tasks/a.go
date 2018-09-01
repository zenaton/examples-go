package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var TaskA = task.New2(&A{})

type A struct{}

func (a *A) Handle() {
	time.Sleep(4 * time.Second)
	fmt.Println("Task A")
}
