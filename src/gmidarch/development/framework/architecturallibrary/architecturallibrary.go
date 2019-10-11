package architecturallibrary

import (
	"errors"
	"gmidarch/development/framework/components"
	"gmidarch/development/framework/connectors"
	"newsolution/shared/parameters"
)

type Record struct {
	Go        interface{}
	Behaviour string
	LTS       interface{}
	PRISM     interface{}
	RBD       interface{} // TODO
	PetriNet  interface{} // TODO
}

type ArchitecturalLibrary struct {
	Lib map[string]Record
}

func (r *Record) SetCSP(csp string) {
	r.CSP = csp
}

func (l *ArchitecturalLibrary) Load() error {
	r1 := *new(error)

	l.Lib = make(map[string]Record)

	// load
	l.Lib["Sender"] = Record{RBD: "TODO", PRISM: "TODO", Go: components.Sender{}, CSP: "B = I_PreInvR1 -> InvR.e1 -> B [] I_PreInvR2 -> InvR.e1 -> B"}
	l.Lib["Receiver"] = Record{RBD: "TODO", PRISM: "TODO", Go: components.Receiver{}, CSP: "B = InvP.e1 -> I_PosInvP -> B"}
	l.Lib["ExecutionEnvironment"] = Record{RBD: "TODO", PRISM: "TODO", Go: components.ExecutionEnvironment{}, CSP: "B = " + parameters.RUNTIME_BEHAVIOUR}
	l.Lib["MAPEKMonitorEvolutive"] = Record{RBD: "TODO", PRISM: "TODO", Go: components.MAPEKMonitorEvolutive{}, CSP: "B = I_EvolutiveMonitoring -> (I_HasPlugin -> InvR.e1 -> B [] I_HasNotPlugin -> B)"}
	l.Lib["MAPEKMonitor"] = Record{RBD: "TODO", PRISM: "TODO", Go: components.MAPEKMonitor{}, CSP: "B = InvP.e1 -> I_Monitor -> InvR.e2 -> B"}
	l.Lib["MAPEKAnalyser"] = Record{RBD: "TODO", PRISM: "TODO", Go: components.MAPEKAnalyser{}, CSP: "B = InvP.e1 -> I_Analyse -> InvR.e2 -> B"}
	l.Lib["MAPEKPlanner"] = Record{RBD: "TODO", PRISM: "TODO", Go: components.MAPEKPlanner{}, CSP: "B = InvP.e1 -> I_Plan -> InvR.e2 -> B"}
	l.Lib["MAPEKExecutor"] = Record{RBD: "TODO", PRISM: "TODO", Go: components.MAPEKExecutor{}, CSP: "B = InvP.e1 -> I_Execute -> InvR.e2 -> B"}
	l.Lib["ExecutionUnit"] = Record{RBD: "TODO", PRISM: "TODO", Go: components.ExecutionUnit{}, CSP: "B = " + parameters.RUNTIME_BEHAVIOUR}
	l.Lib["RequestReply"] = Record{RBD: "TODO", PRISM: "TODO", Go: connectors.RequestReply{}, CSP: "B = InvP.e1 -> InvR.e2 -> TerR.e2 -> TerP.e1 -> B"}
	l.Lib["TwoToOne"] = Record{RBD: "TODO", PRISM: "TODO", Go: connectors.TwoToOne{}, CSP: "B = InvP.e1 -> InvR.e2 -> TerR.e2 -> TerP.e1 -> B [] InvP.e3 -> InvR.e2 -> TerR.e2 -> TerP.e3 -> B"}
	l.Lib["ThreeToOne"] = Record{RBD: "TODO", PRISM: "TODO", Go: connectors.ThreeToOne{}, CSP: "B = InvP.e1 -> InvR.e2 -> B [] InvP.e3 -> InvR.e2 -> B [] InvP.e4 -> InvR.e2 -> B"}
	l.Lib["NTo1"] = Record{RBD: "TODO", PRISM: "TODO", Go: connectors.NTo1{}, CSP: "B = InvP.e1 -> InvR.e2 -> B [] InvP.e3 -> InvR.e2 -> B [] InvP.e4 -> InvR.e2 -> B"}
	l.Lib["OneWay"] = Record{RBD: "TODO", PRISM: "TODO", Go: connectors.OneWay{}, CSP: "B = InvP.e1 -> InvR.e2 -> B"}
	l.Lib["OneToN"] = Record{RBD: "TODO", PRISM: "TODO", Go: connectors.OneToN{}, CSP: "B = " + parameters.RUNTIME_BEHAVIOUR}

	// check
	err := l.CheckLibrary()
	if (err != nil) {
		r1 = errors.New("Architectural Library:: " + err.Error())
		return r1
	}

	return r1
}

func (l ArchitecturalLibrary) CheckLibrary() error {
	r1 := *new(error)

	for elem := range l.Lib {
		if l.Lib[elem].CSP == "" {
			r1 := errors.New("Behaviour '" + elem + "' is INVALID!!")
			return r1
		}
	}
	return r1
}

func (l ArchitecturalLibrary) GetRecord(t string) (error, Record) {
	r1 := *new(error)
	r2 := Record{}

	r2, ok := l.Lib[t]
	if !ok {
		r1 = errors.New("Element type '" + t + "' is NOT in Architectural Library!!")
		return r1, r2
	}
	return r1, r2
}
