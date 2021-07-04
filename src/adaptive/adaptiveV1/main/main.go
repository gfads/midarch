package main

import (
	evolutive2 "adaptive/adaptiveV1/monitor"
	"fmt"
	"sync"
	"time"
)

func Behaviour00() {
	fmt.Println("Hard-coded 00")
}
func Behaviour01() {
	fmt.Println("Hard-coded 01")
}
func Behaviour02() {
	fmt.Println("Hard-coded 02")
}

func AdaptationSystem(in chan int, out chan func(), wg *sync.WaitGroup) { // ADAPTACAO
	for {
		behaviours := []func(){Behaviour00, Behaviour01, Behaviour02}
		tempBehaviour := evolutive2.Monitor{}.LoadFuncs()
		for i := 0; i < len(tempBehaviour); i ++{
			behaviours = append(behaviours, tempBehaviour[i])
		}

		n := <-in
		if n < len(behaviours) {
			out <- behaviours[n]
		} else {
			out <- behaviours[0]
		}
	}
	wg.Done()
}

func ManagedElement(in chan func(), wg *sync.WaitGroup) { // NEGóCIO
	for {
		behaviour := <-in
		fmt.Print("Managed Element: ")
		behaviour()
	}
	wg.Done()
}

func Environment(out chan int, wg *sync.WaitGroup) {
	var n int

	for {
		fmt.Scanln(&n)
		//n = rand.Intn(4)
		out <- n

		time.Sleep(5 * time.Second)
	}
	wg.Done()
}

func main() {
	var wg sync.WaitGroup

	in := make(chan int)
	out := make(chan func())

	wg.Add(3)
	go ManagedElement(out, &wg)
	go AdaptationSystem(in, out, &wg)
	go Environment(in, &wg)

	wg.Wait()
}

