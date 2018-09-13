package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var TaskD = task.NewDefault("TaskD",
	func() (interface{}, error) {

		fmt.Println("Task D starts")
		time.Sleep(9 * time.Second)
		fmt.Println("Task D ends")

		return 3, nil
	})
