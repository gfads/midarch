package main

import (
	"fmt"
	"reflect"
)

type X struct {
	F1 int
}

func main() {
	x := X{F1:1}

	aux := reflect.ValueOf(x).FieldByName("F1")

	fmt.Println(aux)
}
