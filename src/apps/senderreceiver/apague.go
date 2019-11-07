package main

import (
	"fmt"
	"os"
	"reflect"
)

type Elem struct{}

func (Elem) X(){
	fmt.Printf("Here")
}

func main(){
	f,ok := reflect.TypeOf(Elem{}).MethodByName("X")
	if !ok {
		fmt.Printf("Function does not exist!!")
		os.Exit(0)
	}

	fmt.Printf("%v\n",f)

	in := make([]reflect.Value,1)
	in[0] = reflect.ValueOf(Elem{})

	f.Func.Call(in)
}
