package main

import (
	_ "github.com/zenaton/examples-go/client" // initialize zenaton client with credentials
	"github.com/zenaton/examples-go/workflows"
)

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
