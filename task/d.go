package task

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton"
)

var D = zenaton.NewTask(zenaton.TaskParams{
	Name: "TaskD",
	HandleFunc: func() (string, error) {
		fmt.Println("Task D")
		time.Sleep(9 * time.Second)
		return "Task C", nil
	},
})
