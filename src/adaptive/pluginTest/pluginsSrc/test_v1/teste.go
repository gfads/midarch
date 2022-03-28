package test

import "fmt"

type Test struct {
	Name string
	Info string
}

func (t Test) Print() {
	fmt.Println("Here is the Test object:", t.Name, t.Info)
}