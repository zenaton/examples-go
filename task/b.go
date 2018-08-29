package task

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton"
)

var B = zenaton.NewTask(zenaton.TaskParams{
	Name: "TaskB",
	HandleFunc: func() (string, error) {
		time.Sleep(3 * time.Second)
		fmt.Println("Task B")
		return "Task B", nil
	},
})
