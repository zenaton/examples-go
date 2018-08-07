package task

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton"
)

var A = zenaton.NewTask(zenaton.TaskParams{
	Name: "TaskA",
	HandleFunc: func() (string, error) {
		fmt.Println("Task A")
		time.Sleep(3 * time.Second)
		return "Task A", nil
	},
})
