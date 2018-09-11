package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var TaskB = task.NewDefault("TaskB",
	func() (interface{}, error) {
		time.Sleep(3 * time.Second)
		fmt.Println("Task B")
		return nil, nil
	})
