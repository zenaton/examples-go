package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var TaskA = task.NewDefault("TaskA",
	func() (interface{}, error) {
		time.Sleep(4 * time.Second)
		fmt.Println("Task A")
		return nil, nil
	})
