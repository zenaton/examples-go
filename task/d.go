package task

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton"
)

var D = zenaton.NewTask(zenaton.TaskParams{
	Name: "TaskD",
	HandleFunc: func() (string, error) {
		time.Sleep(1 * time.Second)
		fmt.Println("Task D")
		return "Task C", nil
	},
})
