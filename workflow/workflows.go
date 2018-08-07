package workflow

//
//import (
//	"fmt"
//
//	"github.com/zenaton/examples-go/task"
//	"github.com/zenaton/examples-go/tasks"
//	"github.com/zenaton/zenaton-go/v1/zenaton"
//)
//
//var (
//	WaitWorkflow = zenaton.NewWorkflow("WaitWorkflow", func() interface{} {
//		// todo: figure out how to do something like this.email in javascript example
//		task.TaskA.Execute()
//		// todo: kind of ugly to pass in nil here
//		zenaton.Wait(nil).Seconds(5).Execute()
//		task.TaskB.Execute()
//		return nil
//	})
//
//	WaitEventWorkflow = zenaton.NewWorkflow("WaitEventWorkflow", func() interface{} {
//
//		// Wait until the event or 4 seconds
//		event, err := zenaton.Wait("MyEvent").Seconds(4).Execute()
//		if err != nil {
//			panic(err)
//		}
//
//		// If event has been triggered
//		if event != nil {
//			// Execute TaskB
//			task.TaskA.Execute()
//		} else {
//			// Execute Task B
//			task.TaskB.Execute()
//		}
//		return nil
//	}).IDFunc(func() string {
//		return "MyId"
//	})
//)
