package workflow

import (
	"fmt"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var DisplayTask = task.NewCustom("DisplayTask", &Display{})

type Display struct {
	ID int
}

func (dt *Display) Init(id int) {
	dt.ID = id
}

func (dt *Display) Handle() (interface{}, error) {
	fmt.Print(dt.ID)
	return nil, nil
}
