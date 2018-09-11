package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var TaskC = task.NewDefault("TaskC",
	func() (interface{}, error) {
		time.Sleep(2 * time.Second)
		fmt.Println("Task C")
		return nil, nil
	})
