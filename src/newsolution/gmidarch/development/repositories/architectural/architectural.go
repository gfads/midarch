package architectural

import (
	"newsolution/gmidarch/development/components"
	"newsolution/gmidarch/development/connectors"
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
	al.Library["Oneway"] = Record{Type: connectors.NewOneway(), Behaviour: connectors.NewOneway().Behaviour}
	al.Library["Requestreply"] = Record{Type: connectors.NewRequestReply(), Behaviour: connectors.NewRequestReply().Behaviour}
	al.Library["Sender"] = Record{Type: components.NewSender(), Behaviour: components.NewSender().Behaviour}
	al.Library["Receiver"] = Record{Type: components.NewReceiver(), Behaviour: components.NewReceiver().Behaviour}
	al.Library["Client"] = Record{Type: components.NewClient(), Behaviour: components.NewClient().Behaviour}
	al.Library["Server"] = Record{Type: components.NewServer(), Behaviour: components.NewServer().Behaviour}
	al.Library["Calculatorproxy"] = Record{Type: components.NewCalculatorProxy(), Behaviour: ""}
	al.Library["Marshaller"] = Record{Type: components.NewMarshaller(), Behaviour: ""}
	al.Library["Requestor"] = Record{Type: components.NewRequestor(), Behaviour: ""}
	al.Library["CRH"] = Record{Type: components.NewCRH(), Behaviour: components.NewCRH().Behaviour}
	al.Library["SRH"] = Record{Type: components.NewSRH(), Behaviour: components.NewSRH().Behaviour}
	al.Library["Calculatorserver"] = Record{Type: components.Newcalculatorserver(), Behaviour: ""}
	al.Library["Calculatorinvoker"] = Record{Type: components.NewCalculatorinvoker(), Behaviour: ""}
	al.Library["Calculatorclient"] = Record{Type: components.NewCalculatorclient(), Behaviour: components.NewCalculatorclient().Behaviour}
	al.Library["Core"] = Record{Type: components.NewCore(), Behaviour: components.NewCore().Behaviour}
	al.Library["ComponenteNovo"] = Record{Type: components.NewCore(), Behaviour: components.NewComponentNovo().Behaviour}

	return r1
}
