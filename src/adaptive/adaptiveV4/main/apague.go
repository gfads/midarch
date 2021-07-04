package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Person struct {
	Name string
	Age int
}

func main() {
	p := Person{Name:"Jose",Age:25}

	fmt.Println(reflect.TypeOf(p))

	fmt.Print("Enter a Method Name: ")
	m := MyRead()

	Invoke(p,m,)
}

func (p Person) SayHello(){
	fmt.Println("Hello")
}
func (p Person) SayHi() {
	fmt.Println("Hi!!")
}

func MyRead() string {
	reader := bufio.NewReader(os.Stdin)
	d, _ := reader.ReadString('\n')
	x := strings.Split(d,"\n")
	m := x[0]

	return m
}

func Invoke(any interface{}, name string, args... interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	reflect.ValueOf(any).MethodByName(name).Call(inputs)
}


