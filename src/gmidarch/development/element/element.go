package element

import (
	messages2 "gmidarch/development/messages"
)

type Element struct{}

func (Element) InvR(invR *chan messages2.SAMessage, msg *messages2.SAMessage) {
	*invR <- *msg
}

func (Element) TerR(terR *chan messages2.SAMessage, msg *messages2.SAMessage) {
	*msg = <-*terR
}

func (Element) InvP(invP *chan messages2.SAMessage, msg *messages2.SAMessage) {
	*msg = <-*invP
}

func (Element) TerP(terP *chan messages2.SAMessage, msg *messages2.SAMessage) {
	*terP <- *msg
}
