package architectural

import (
	components2 "gmidarch/development/components"
	connectors2 "gmidarch/development/connectors"
)

type Record struct {
	Type      interface{}
	Behaviour string
}

type ArchitecturalRepository struct {
	Library map[string]Record
}

func (al *ArchitecturalRepository) Load() error {
	r1 := *new(error)

	al.Library = make(map[string]Record)

	// load
	al.Library["Analyser"] = Record{Type: components2.NewAnalyser(), Behaviour: components2.NewAnalyser().Behaviour}
	al.Library["OnetoN"] = Record{Type: connectors2.NewOnetoN(), Behaviour: connectors2.NewOnetoN().Behaviour}
	al.Library["Oneway"] = Record{Type: connectors2.NewOneway(), Behaviour: connectors2.NewOneway().Behaviour}
	al.Library["Receiver"] = Record{Type: components2.NewReceiver(), Behaviour: components2.NewReceiver().Behaviour}
	al.Library["Requestreply"] = Record{Type: connectors2.NewRequestReply(), Behaviour: connectors2.NewRequestReply().Behaviour}
	al.Library["Sender"] = Record{Type: components2.NewSender(), Behaviour: components2.NewSender().Behaviour}
	al.Library["Client"] = Record{Type: components2.NewClient(), Behaviour: components2.NewClient().Behaviour}
	al.Library["Server"] = Record{Type: components2.NewServer(), Behaviour: components2.NewServer().Behaviour}
	al.Library["Calculatorproxy"] = Record{Type: components2.NewCalculatorProxy(), Behaviour: components2.NewCalculatorProxy().Behaviour}
	al.Library["Marshaller"] = Record{Type: components2.NewMarshaller(), Behaviour: components2.NewMarshaller().Behaviour}
	al.Library["Requestor"] = Record{Type: components2.NewRequestor(), Behaviour: components2.NewRequestor().Behaviour}
	al.Library["CRH"] = Record{Type: components2.NewCRH(), Behaviour: components2.NewCRH().Behaviour}
	al.Library["SRH"] = Record{Type: components2.NewSRH(), Behaviour: components2.NewSRH().Behaviour}
	al.Library["Calculatorserver"] = Record{Type: components2.Newcalculatorserver(), Behaviour: components2.Newcalculatorserver().Behaviour}
	al.Library["Calculatorinvoker"] = Record{Type: components2.NewCalculatorinvoker(), Behaviour: components2.NewCalculatorinvoker().Behaviour}
	al.Library["Calculatorclient"] = Record{Type: components2.NewCalculatorclient(), Behaviour: components2.NewCalculatorclient().Behaviour}
	al.Library["Core"] = Record{Type: components2.NewCore(), Behaviour: components2.NewCore().Behaviour}
	al.Library["Unit"] = Record{Type: components2.NewUnit(), Behaviour: components2.NewUnit().Behaviour}
	al.Library["Monevolutive"] = Record{Type: components2.NewMonevolutive(), Behaviour: components2.NewMonevolutive().Behaviour}
	al.Library["Monitor"] = Record{Type: components2.NewMonitor(), Behaviour: components2.NewMonitor().Behaviour}
	al.Library["Planner"] = Record{Type: components2.NewPlanner(), Behaviour: components2.NewPlanner().Behaviour}
	al.Library["Executor"] = Record{Type: components2.NewExecutor(), Behaviour: components2.NewExecutor().Behaviour}

	return r1
}
