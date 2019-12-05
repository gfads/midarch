package components

import (
	"fmt"
	"gmidarch/development/connectors"
	"gmidarch/development/messages"
	"os"
)

func ConfigureSelector(c string) func (interface{}, [] *interface{},string, *messages.SAMessage,[]*interface{}) {
    s := func (interface{}, [] *interface{},string, *messages.SAMessage,[]*interface{}){}

		switch c {
		case "Sender":
			s = NewSender().Selector
		case "Receiver":
			s = NewReceiver().Selector
		case "Unit":
			s = NewUnit().Selector
		case "Core":
			s = NewCore().Selector
		case "Client":
			s = NewClient().Selector
		case "Server":
			s = NewServer().Selector
		case "Fibonacciserver":
			s = Newfibonacciserver().Selector
		case "FibonacciinvokerM":
			s = NewFibonacciInvokerM().Selector
		case "SRH":
			s = NewSRH().Selector
		case "Fibonacciclient":
			s = NewFibonacciclient().Selector
		case "Fibonacciproxy":
			s = NewFibonacciproxy().Selector
		case "RequestorM":
			s = NewRequestorM().Selector
		case "CRH":
			s = NewCRH().Selector
		case "Namingserver":
			s = Newnamingserver().Selector
		case "NaminginvokerM":
			s = NewnaminginvokerM().Selector
		case "Namingproxy":
			s = NewNamingproxy().Selector
		case "Monevolutive":
			s = NewMonevolutive().Selector
		case "Monitor":
			s = NewMonitor().Selector
		case "Analyser":
			s = NewAnalyser().Selector
		case "Planner":
			s = NewPlanner().Selector
		case "Executor":
			s = NewExecutor().Selector
		case "Oneto8":
			s = connectors.NewOneto8().Selector
		case "Oneto5":
			s = connectors.NewOneto5().Selector
		case "Oneto3":
			s = connectors.NewOneto3().Selector
		default:
			fmt.Printf("Generic:: Element '%v' has not a selector!!\n", c)
			os.Exit(0)
		}

	return s
}
