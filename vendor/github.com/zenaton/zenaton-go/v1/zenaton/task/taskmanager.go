package task

import (
	"fmt"
	"github.com/zenaton/zenaton-go/v1/zenaton/service/serializer"
	"sync"
)

// UnsafeManager is used by the agent, and thus must be exported. But a normal user of the library shouldn't use this
// directly.
var UnsafeManager = &Store{
	tasks: make(map[string]*Definition),
	mu:    &sync.RWMutex{},
}

// Store is a thread-safe store of task Definitions. This is used to insure that no two tasks can have the same name.
// It also will be used by the agent to be able take a task name (as well as any task data if it exists) and produce an
// Instance of that task.
type Store struct {
	tasks map[string]*Definition
	mu    *sync.RWMutex
}

func (s *Store) setDefinition(name string, tt *Definition) {
	// check that this task does not exist yet
	if s.UnsafeGetDefinition(name) != nil {
		panic(fmt.Sprint("Instance definition with name '", name, "' already exists"))
	}

	s.mu.Lock()
	s.tasks[name] = tt
	s.mu.Unlock()

}

// UnsafeGetDefinition is used by the agent, and thus must be exported. But a normal user of the library shouldn't use this
// directly. UnsafeGetDefinition takes a name, and retrieves the corresponding task Definition from the store.
func (s *Store) UnsafeGetDefinition(name string) *Definition {
	s.mu.RLock()
	t := s.tasks[name]
	s.mu.RUnlock()
	return t
}

// UnsafeGetInstance is used by the agent, and thus must be exported. But a normal user of the library shouldn't use this
// directly. UnsafeGetInstance takes a name of a task, and the task's data, and can create an Instance the task.
func (s *Store) UnsafeGetInstance(name, encodedData string) *Instance {

	// get task class
	tt := s.UnsafeGetDefinition(name)

	// unserialize data
	err := serializer.Decode(encodedData, tt.defaultTask.Handler)
	if err != nil {
		panic(err)
	}

	return tt.defaultTask
}
