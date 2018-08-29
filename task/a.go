package task

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton"
)

var A = zenaton.NewTask(zenaton.TaskParams{
	Name: "TaskA",
	HandleFunc: func() (string, error) {
		time.Sleep(4 * time.Second)
		fmt.Println("Task A")
		return "Task A", nil
	},
})
