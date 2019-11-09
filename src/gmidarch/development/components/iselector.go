package components

import "gmidarch/development/messages"

type Selector interface {
	Selector(elem interface{}, op string) func(*messages.SAMessage, []*interface{})
}
