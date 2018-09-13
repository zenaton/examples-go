package task

import (
	"fmt"
	"github.com/zenaton/zenaton-go/v1/zenaton/service/serializer"
	"sync"
)

var taskManagerInstance = &TaskManager{
	tasks: make(map[string]*taskType),
	mu:    &sync.RWMutex{},
}

type TaskManager struct {
	tasks map[string]*taskType
	mu    *sync.RWMutex
}

//todo: problem, This shouldn't be accessible to the user
func NewTaskManager() *TaskManager {
	return taskManagerInstance
}

func (tm *TaskManager) setClass(name string, tt *taskType) {
	// check that this task does not exist yet
	//todo: is this right?
	if tm.GetClass(name) != nil {
		panic(fmt.Sprint("Task definition with name '", name, "' already exists"))
	}

	tm.mu.Lock()
	tm.tasks[name] = tt
	tm.mu.Unlock()

}

func (tm *TaskManager) GetClass(name string) *taskType {
	tm.mu.RLock()
	t := tm.tasks[name]
	tm.mu.RUnlock()
	return t
}

func (tm *TaskManager) GetTask(name, encodedData string) *Task {

	// get task class
	tt := tm.GetClass(name)

	// unserialize data
	err := serializer.Decode(encodedData, tt.defaultTask.Handler)
	if err != nil {
		panic(err)
	}

	return tt.defaultTask
}
