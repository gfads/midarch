package main

import (
	"fmt"
	"testing"
)

func TestHelloWorld(t *testing.T){
	for i := 0; i<10000; i++{
		fmt.Printf("Nothing to do\n")
	}
}
