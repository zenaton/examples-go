package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var TaskA = task.NewDefault("TaskA",

	func() (interface{}, error) {

		fmt.Println("Task A starts")
		time.Sleep(3 * time.Second)
		fmt.Println("Task A ends")

		return 0, nil
	})
