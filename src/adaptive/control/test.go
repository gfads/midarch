package main

import (
	"fmt"
	"math/rand"
)

func g(s int, ch chan int) {
	n := 0
	for {
		r := <-ch
		if r == s {
			fmt.Println("Go routine", s, n)
			n++
		}
	}
}

func scheduler(ch chan int) {
	for {
		ch <- rand.Intn(3)
	}
}

func main() {
	ch := make(chan int)

	//runtime.GOMAXPROCS(2)
	go g(1, ch)
	go g(2, ch)
	go scheduler(ch)

	fmt.Scanln()
}
