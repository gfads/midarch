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

func (u Element) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {

	switch reflect.TypeOf(elem).Name() {
	case "Sender":
		components.Sender{}.Selector(elem, elemInfo, op, msg, info, r)
	case "Receiver":
		components.Receiver{}.Selector(elem, elemInfo, op, msg, info, r)
	case "Unit":
		components.Unit{}.Selector(elem, elemInfo, op, msg, info, r)
	case "Core":
		components.Core{}.Selector(elem, elemInfo, op, msg, info, r)
	case "Client":
		components.Client{}.Selector(elem, elemInfo, op, msg, info, r)
	case "Server":
		components.Server{}.Selector(elem, elemInfo, op, msg, info, r)
	case "Fibonacciserver":
		components.Fibonacciserver{}.Selector(elem, elemInfo, op, msg, info, r)
	case "FibonacciinvokerM":
		components.FibonacciinvokerM{}.Selector(elem, elemInfo, op, msg, info, r)
	case "SRH":
		components.SRH{}.Selector(elem, elemInfo, op, msg, info, r)
	case "Fibonacciclient":
		components.Fibonacciclient{}.Selector(elem, elemInfo, op, msg, info, r)
	case "Fibonacciproxy":
		components.Fibonacciproxy{}.Selector(elem, elemInfo, op, msg, info, r)
	case "RequestorM":
		components.RequestorM{}.Selector(elem, elemInfo, op, msg, info, r)
	case "CRH":
		components.CRH{}.Selector(elem, elemInfo, op, msg, info, r)
	case "Namingserver":
		components.Namingserver{}.Selector(elem, elemInfo, op, msg, info, r)
	case "NaminginvokerM":
		components.NaminginvokerM{}.Selector(elem, elemInfo, op, msg, info, r)
	case "Namingproxy":
		components.Namingproxy{}.Selector(elem, elemInfo, op, msg, info, r)
	case "Monevolutive":
		components.Monevolutive{}.Selector(elem, elemInfo, op, msg, info, r)
	case "Monitor":
		components.Monitor{}.Selector(elem, elemInfo, op, msg, info, r)
	case "Analyser":
		components.Analyser{}.Selector(elem, elemInfo, op, msg, info, r)
	case "Planner":
		components.Planner{}.Selector(elem, elemInfo, op, msg, info, r)
	case "Executor":
		components.Executor{}.Selector(elem, elemInfo, op, msg, info, r)
	default:
		fmt.Printf("Exec:: Element '%v' not visible in Engine!!\n", reflect.TypeOf(elem).Name())
		os.Exit(0)
	}
}
