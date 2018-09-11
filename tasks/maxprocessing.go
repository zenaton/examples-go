package tasks

import (
	"fmt"
	"time"

	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var MaxProcessingTask = task.New("MaxProcessingTask", &MaxProcessing{})

type MaxProcessing struct {
	Max int64
}

func (mp *MaxProcessing) Handle() (interface{}, error) {
	time.Sleep(time.Duration(mp.Max+6) * time.Second)
	fmt.Println("This Shouldn't Print")
	return nil, nil
}

func (mp *MaxProcessing) MaxTime() int64 {
	return mp.Max
}
