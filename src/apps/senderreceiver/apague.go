package main

import (
	"fmt"
	"os"
	"reflect"
	"time"
)

type I1 interface{
	Selector(int) func()
}

type Elem1 struct{}

type Elem2 struct{}

func (Elem1) Selector(n int) func(){

	var f func()

	switch n {
	case 1:
		f = func(){
			Elem1{}.X()
		}
	}
	return f
}

func (Elem2) Selector(n int) func(){

	var f func()

	switch n {
	case 1:
		f = func(){
			Elem2{}.X()
		}
	}
	return f
}

func (Elem1) X() {
	fmt.Printf("X\n")
}

func (Elem2) X() {
	fmt.Printf("X\n")
}

func main() {
	f, ok := reflect.TypeOf(Elem1{}).MethodByName("X")
	if !ok {
		fmt.Printf("Function does not exist!!")
		os.Exit(0)
	}

	fmt.Printf("%v\n", f)

	in := make([]reflect.Value, 1)
	in[0] = reflect.ValueOf(Elem1{})

	t1 := time.Now()
	var s I1
	for i := 0; i < 100000; i++ {
		//f.Func.Call(in)
		//Elem{}.X()
		s = Elem2{}
		f := s.Selector(1)
		f()
	}
	fmt.Printf("Total Time: %v\n", time.Now().Sub(t1))

	// elemType := tp.(func()interface{})()

}
