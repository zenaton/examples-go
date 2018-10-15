package tasks

import (
	"fmt"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var E = task.New("TaskE",
	func() (interface{}, error) {

		fmt.Println("Task E starts")
		panic("Error in task E")
		fmt.Println("Task E ends")

		return nil, nil
	})
