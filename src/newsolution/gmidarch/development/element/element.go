package element

import (
	"newsolution/gmidarch/development/messages"
)

type Element struct{}

func (Element) InvR(invR *chan messages.SAMessage, msg *messages.SAMessage) {
	*invR <- *msg
}

func (Element) TerR(terR *chan messages.SAMessage, msg *messages.SAMessage) {
	*msg = <-*terR
}

func (Element) InvP(invP *chan messages.SAMessage, msg *messages.SAMessage) {
	*msg = <-*invP
}

func (Element) TerP(terP *chan messages.SAMessage, msg *messages.SAMessage) {
	*terP <- *msg
}
