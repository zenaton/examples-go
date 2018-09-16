package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var C = task.New("TaskC",

	func() (interface{}, error) {

		fmt.Println("Task C starts")
		time.Sleep(7 * time.Second)
		fmt.Println("Task C ends")

		return 2, nil
	})
