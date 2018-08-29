package task

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton"
)

var C = zenaton.NewTask(zenaton.TaskParams{
	Name: "TaskC",
	HandleFunc: func() (string, error) {
		time.Sleep(2 * time.Second)
		fmt.Println("Task C")
		return "Task C", nil
	},
})
