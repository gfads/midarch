package components

import "gmidarch/development/messages"

type Selector interface {
	Selector(interface{}, [] *interface{}, string, *messages.SAMessage, []*interface{})
}

func NewSelector() Selector {
	r := new(Selector)

	return *r
}