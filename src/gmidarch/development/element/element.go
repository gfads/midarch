package element

import (
	"fmt"
	"gmidarch/development/components"
	messages "gmidarch/development/messages"
	"os"
	"reflect"
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

func (u Element) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}) {

	switch reflect.TypeOf(elem).Name() {
	case "Sender":
		components.Sender{}.Selector(elem, elemInfo, op, msg, info)
	case "Receiver":
		components.Receiver{}.Selector(elem, elemInfo, op, msg, info)
	case "Unit":
		components.Unit{}.Selector(elem, elemInfo, op, msg, info)
	case "Core":
		components.Core{}.Selector(elem, elemInfo, op, msg, info)
	case "Client":
		components.Client{}.Selector(elem, elemInfo, op, msg, info)
	case "Server":
		components.Server{}.Selector(elem, elemInfo, op, msg, info)
	case "Fibonacciserver":
		components.Fibonacciserver{}.Selector(elem, elemInfo, op, msg, info)
	case "FibonacciinvokerM":
		components.FibonacciinvokerM{}.Selector(elem, elemInfo, op, msg, info)
	case "SRH":
		components.SRH{}.Selector(elem, elemInfo, op, msg, info)
	case "Fibonacciclient":
		components.Fibonacciclient{}.Selector(elem, elemInfo, op, msg, info)
	case "Fibonacciproxy":
		components.Fibonacciproxy{}.Selector(elem, elemInfo, op, msg, info)
	case "RequestorM":
		components.RequestorM{}.Selector(elem, elemInfo, op, msg, info)
	case "CRH":
		components.CRH{}.Selector(elem, elemInfo, op, msg, info)
	case "Namingserver":
		components.Namingserver{}.Selector(elem, elemInfo, op, msg, info)
	case "NaminginvokerM":
		components.NaminginvokerM{}.Selector(elem, elemInfo, op, msg, info)
	case "Namingproxy":
		components.Namingproxy{}.Selector(elem, elemInfo, op, msg, info)
	case "Monevolutive":
		components.Monevolutive{}.Selector(elem, elemInfo, op, msg, info)
	case "Monitor":
		components.Monitor{}.Selector(elem, elemInfo, op, msg, info)
	case "Analyser":
		components.Analyser{}.Selector(elem, elemInfo, op, msg, info)
	case "Planner":
		components.Planner{}.Selector(elem, elemInfo, op, msg, info)
	case "Executor":
		components.Executor{}.Selector(elem, elemInfo, op, msg, info)
	default:
		fmt.Printf("Exec:: Element '%v' not visible in Engine!!\n", reflect.TypeOf(elem).Name())
		os.Exit(0)
	}
}
