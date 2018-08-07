package task

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton"
)

var C = zenaton.NewTask(zenaton.TaskParams{
	Name: "TaskC",
	HandleFunc: func() (string, error) {
		fmt.Println("Task C")
		time.Sleep(7 * time.Second)
		return "Task C", nil
	},
})
