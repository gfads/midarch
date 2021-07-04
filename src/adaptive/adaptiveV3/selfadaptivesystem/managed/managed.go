package managed

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Managed struct{}

func (m Managed) Start(toManaging chan int, fromManaging chan func(), wg *sync.WaitGroup) {
	f := func() {
		fmt.Println("Default behaviour")
	}
	for {
		select {
		case p := <-fromManaging:
			f = p
		//case toManaging <- m.sense():  TODO
		default:
		}
		f()
		time.Sleep(500 * time.Millisecond)
	}
}

func (Managed) sense() int {
	r := rand.Intn(10) // TODO
	return r
}
