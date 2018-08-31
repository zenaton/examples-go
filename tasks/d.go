package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

func init() {
	//todo: call this registerTask
	task.New(&D{})
}

type D struct{}

func (a *D) Handle() {
	time.Sleep(1 * time.Second)
	fmt.Println("Task D")
}
