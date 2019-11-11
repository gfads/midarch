package components

import "gmidarch/development/messages"

type Selector interface {
	Selector(interface{}, string, *messages.SAMessage, []*interface{})
}
