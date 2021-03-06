package architectural

import (
	"gmidarch/development/components"
	"gmidarch/development/connectors"
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
	al.Library["Analyser"] = Record{Type: components.NewAnalyser(), Behaviour: components.NewAnalyser().Behaviour}
	al.Library["OnetoN"] = Record{Type: connectors.NewOnetoN(), Behaviour: connectors.NewOnetoN().Behaviour}
	al.Library["Ntoone"] = Record{Type: connectors.NewNtoOne(), Behaviour: connectors.NewNtoOne().Behaviour}
	al.Library["Oneto2"] = Record{Type: connectors.NewOneto2(), Behaviour: connectors.NewOneto2().Behaviour}
	al.Library["Oneto8"] = Record{Type: connectors.NewOneto8(), Behaviour: connectors.NewOneto8().Behaviour}
	al.Library["Oneto5"] = Record{Type: connectors.NewOneto5(), Behaviour: connectors.NewOneto5().Behaviour}
	al.Library["Oneto3"] = Record{Type: connectors.NewOneto3(), Behaviour: connectors.NewOneto3().Behaviour}
	al.Library["Oneto6"] = Record{Type: connectors.NewOneto6(), Behaviour: connectors.NewOneto6().Behaviour}
	al.Library["Oneto7"] = Record{Type: connectors.NewOneto7(), Behaviour: connectors.NewOneto7().Behaviour}
	al.Library["Oneto9"] = Record{Type: connectors.NewOneto9(), Behaviour: connectors.NewOneto9().Behaviour}
	al.Library["Oneway"] = Record{Type: connectors.NewOneway(), Behaviour: connectors.NewOneway().Behaviour}
	al.Library["Receiver"] = Record{Type: components.NewReceiver(), Behaviour: components.NewReceiver().Behaviour}
	al.Library["Requestreply"] = Record{Type: connectors.Newrequestreply(), Behaviour: connectors.Newrequestreply().Behaviour}
	al.Library["Reqreponeto2"] = Record{Type: connectors.Newreqreponeto2(), Behaviour: connectors.Newreqreponeto2().Behaviour}
	al.Library["Sender"] = Record{Type: components.NewSender(), Behaviour: components.NewSender().Behaviour}
	al.Library["Client"] = Record{Type: components.NewClient(), Behaviour: components.NewClient().Behaviour}
	al.Library["Server"] = Record{Type: components.NewServer(), Behaviour: components.NewServer().Behaviour}
	al.Library["Calculatorproxy"] = Record{Type: components.NewCalculatorProxy(), Behaviour: components.NewCalculatorProxy().Behaviour}
	al.Library["Marshaller"] = Record{Type: components.NewMarshaller(), Behaviour: components.NewMarshaller().Behaviour}
	al.Library["Requestor"] = Record{Type: components.NewRequestor(), Behaviour: components.NewRequestor().Behaviour}
	al.Library["RequestorM"] = Record{Type: components.NewRequestorM(), Behaviour: components.NewRequestorM().Behaviour}
	al.Library["CRH"] = Record{Type: components.NewCRH(), Behaviour: components.NewCRH().Behaviour}
	al.Library["SRH"] = Record{Type: components.NewSRH(), Behaviour: components.NewSRH().Behaviour}
	al.Library["CRHHttp"] = Record{Type: components.NewCRHHttp(), Behaviour: components.NewCRHHttp().Behaviour}
	al.Library["SRHHttp"] = Record{Type: components.NewSRHHttp(), Behaviour: components.NewSRHHttp().Behaviour}
	al.Library["CRHHttps"] = Record{Type: components.NewCRHHttps(), Behaviour: components.NewCRHHttps().Behaviour}
	al.Library["CRHHttp2"] = Record{Type: components.NewCRHHttp2(), Behaviour: components.NewCRHHttp2().Behaviour}
	al.Library["SRHHttps"] = Record{Type: components.NewSRHHttps(), Behaviour: components.NewSRHHttps().Behaviour}
	al.Library["SRHHttp2"] = Record{Type: components.NewSRHHttp2(), Behaviour: components.NewSRHHttp2().Behaviour}
	al.Library["CRHSsl"] = Record{Type: components.NewCRHSsl(), Behaviour: components.NewCRHSsl().Behaviour}
	al.Library["SRHSsl"] = Record{Type: components.NewSRHSsl(), Behaviour: components.NewSRHSsl().Behaviour}
	al.Library["CRHQuic"] = Record{Type: components.NewCRHQuic(), Behaviour: components.NewCRHQuic().Behaviour}
	al.Library["SRHQuic"] = Record{Type: components.NewSRHQuic(), Behaviour: components.NewSRHQuic().Behaviour}
	al.Library["CRHUdp"] = Record{Type: components.NewCRHUdp(), Behaviour: components.NewCRHUdp().Behaviour}
	al.Library["SRHUdp"] = Record{Type: components.NewSRHUdp(), Behaviour: components.NewSRHUdp().Behaviour}
	al.Library["CRHRpc"] = Record{Type: components.NewCRHRpc(), Behaviour: components.NewCRHRpc().Behaviour}
	al.Library["SRHRpc"] = Record{Type: components.NewSRHRpc(), Behaviour: components.NewSRHRpc().Behaviour}
	al.Library["RPCRequestorM"] = Record{Type: components.NewRPCRequestorM(), Behaviour: components.NewRPCRequestorM().Behaviour}
	al.Library["RPCInvokerM"] = Record{Type: components.NewRPCInvokerM(), Behaviour: components.NewRPCInvokerM().Behaviour}
	al.Library["HttpRequestorM"] = Record{Type: components.NewHttpRequestorM(), Behaviour: components.NewHttpRequestorM().Behaviour}
	al.Library["Http2RequestorM"] = Record{Type: components.NewHttp2RequestorM(), Behaviour: components.NewHttp2RequestorM().Behaviour}
	al.Library["HttpInvokerM"] = Record{Type: components.NewHttpInvokerM(), Behaviour: components.NewHttpInvokerM().Behaviour}
	al.Library["Http2InvokerM"] = Record{Type: components.NewHttp2InvokerM(), Behaviour: components.NewHttp2InvokerM().Behaviour}
	al.Library["HttpProxy"] = Record{Type: components.NewHttpProxy(), Behaviour: components.NewHttpProxy().Behaviour}
	al.Library["Http2Proxy"] = Record{Type: components.NewHttp2Proxy(), Behaviour: components.NewHttp2Proxy().Behaviour}
	al.Library["Calculatorserver"] = Record{Type: components.Newcalculatorserver(), Behaviour: components.Newcalculatorserver().Behaviour}
	al.Library["Calculatorinvoker"] = Record{Type: components.NewCalculatorinvoker(), Behaviour: components.NewCalculatorinvoker().Behaviour}
	al.Library["Calculatorclient"] = Record{Type: components.NewCalculatorclient(), Behaviour: components.NewCalculatorclient().Behaviour}
	al.Library["Core"] = Record{Type: components.NewCore(), Behaviour: components.NewCore().Behaviour}
	al.Library["Unit"] = Record{Type: components.NewUnit(), Behaviour: components.NewUnit().Behaviour}
	al.Library["Monevolutive"] = Record{Type: components.NewMonevolutive(), Behaviour: components.NewMonevolutive().Behaviour}
	al.Library["Monitor"] = Record{Type: components.NewMonitor(), Behaviour: components.NewMonitor().Behaviour}
	al.Library["Planner"] = Record{Type: components.NewPlanner(), Behaviour: components.NewPlanner().Behaviour}
	al.Library["Executor"] = Record{Type: components.NewExecutor(), Behaviour: components.NewExecutor().Behaviour}
	al.Library["Fibonacciserver"] = Record{Type: components.Newfibonacciserver(), Behaviour: components.Newfibonacciserver().Behaviour}
	al.Library["Fibonacciclient"] = Record{Type: components.NewFibonacciclient(), Behaviour: components.NewFibonacciclient().Behaviour}
	al.Library["Fibonacciinvoker"] = Record{Type: components.NewFibonacciinvoker(), Behaviour: components.NewFibonacciinvoker().Behaviour}
	al.Library["Fibonacciinvokerm"] = Record{Type: components.NewFibonacciInvokerM(), Behaviour: components.NewFibonacciInvokerM().Behaviour}
	al.Library["Fibonacciproxy"] = Record{Type: components.NewFibonacciproxy(), Behaviour: components.NewFibonacciproxy().Behaviour}
	al.Library["Namingserver"] = Record{Type: components.Newnamingserver(), Behaviour: components.Newnamingserver().Behaviour}
	al.Library["Naminginvokerm"] = Record{Type: components.NewnaminginvokerM(), Behaviour: components.NewnaminginvokerM().Behaviour}
	al.Library["Namingproxy"] = Record{Type: components.NewNamingproxy(), Behaviour: components.NewNamingproxy().Behaviour}
	al.Library["Notificationenginex"] = Record{Type: components.NewnotificationengineX(), Behaviour: components.NewnotificationengineX().Behaviour}
	al.Library["Notificationengineinvoker"] = Record{Type: components.Newnotificationengineinvoker(), Behaviour: components.Newnotificationengineinvoker().Behaviour}
	al.Library["Notificationengineproxy"] = Record{Type: components.Newnotificationengineproxy(), Behaviour: components.Newnotificationengineproxy().Behaviour}

	return r1
}
