package managed

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type ManagedInfo struct {
	N             int
}

var chn chan func()

type Managed interface {
	Start(*sync.WaitGroup)
	Sense() ManagedInfo
	Adapt(func())
}

type ManagedImpl struct {
	Behaviour func()
}

func NewManaged() Managed {
	var m Managed

	x := fDefault

	m = &ManagedImpl{Behaviour: x}

	return m
}

func (m *ManagedImpl) Start(wg *sync.WaitGroup) {
	chn = make(chan func())
	for {
		select {
		case m.Behaviour = <-chn:
		default:
		}
		m.Behaviour()
		time.Sleep(1000 * time.Millisecond)
	}
	wg.Done()
}

func (m ManagedImpl) Sense() ManagedInfo {

	r := ManagedInfo{N: rand.Intn(10)} // TODO

	return r
}

func (ManagedImpl) Adapt(f func()) {
	chn <- f
}

func fDefault() {
	fmt.Println("Default")
}
