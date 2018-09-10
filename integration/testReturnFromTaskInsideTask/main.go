package main

import (
	"github.com/zenaton/examples-go/integration/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflows"
)

func init() { client.SetEnv("testReturnFromTaskInsideTask.env") }
func main() {
	//var i interface{}
	//x(&i)
	//fmt.Println("I: ", i)

	workflows.TestRelaunchTaskWorkflow.NewInstance().Dispatch()
}

//
//func x(i interface{}) {
//	var bob interface{}
//	bob = "carl"
//	v := reflect.ValueOf(i)
//	v.Elem().Set(reflect.ValueOf(bob))
//}
