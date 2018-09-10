package log

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var PrintTask = task.New("PrintTask", &Print{})

type Print struct {
	Values string
}

func (p *Print) Handle() (interface{}, error) {
	return fmt.Println(p.Values)
}

func Println(values ...interface{}) error {
	str := spew.Sprint(values...)
	return PrintTask.NewInstance(&Print{str}).Execute()
}
