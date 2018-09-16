package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var B = task.New("TaskB",

	func() (interface{}, error) {

		fmt.Println("Task B starts")
		time.Sleep(5 * time.Second)
		fmt.Println("Task B ends")

		return 1, nil
	})
