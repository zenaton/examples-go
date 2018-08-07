package task

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton"
)

func Binput(input string) *zenaton.Task {
	return zenaton.NewTask(zenaton.TaskParams{
		Data: input,
		Name: "TaskB",
		HandleFunc: func() (string, error) {
			fmt.Println(input + "Task B")
			time.Sleep(5 * time.Second)
			return "Task B", nil
		},
	})
}

var B = zenaton.NewTask(zenaton.TaskParams{
	Name: "TaskB",
	HandleFunc: func() (string, error) {
		fmt.Println("Task B")
		time.Sleep(5 * time.Second)
		return "Task B", nil
	},
})
