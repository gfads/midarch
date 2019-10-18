package main

import (
"fmt"
	"gmidarch/development/framework/messages"
)

type Receiver struct{}

func GetTypeElem() interface{}{
	return Receiver{}
}

func I_APAGUE() string {
	return "RECEIVER 01"
}

func GetBehaviourExp() string {
	//return libraries.BehaviourLibrary[calculatorinvoker.CalculatorInvoker{}]
	return "B = InvP.e1 -> I_PosInvP -> B"
}

func (Receiver) I_PosInvP(msg messages.SAMessage, info interface{}, r *bool) {
	//fmt.Println("Receiver:: " + msg.Payload.(string))
	fmt.Printf("Receiver [TYPE 01]:: %v\n",msg.Payload)
	return
}

