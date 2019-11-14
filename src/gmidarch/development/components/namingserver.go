package components

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"shared"
	"shared/ior"
)

type Namingserver struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func Newnamingserver() Namingserver {

	r := new(Namingserver)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"

	return *r
}

func (e Namingserver) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}) {
	e.I_Process(msg, info)
}

func (e Namingserver) I_Process(msg *messages.SAMessage, info [] *interface{}) {
	request := msg.Payload.(shared.Request)

	switch request.Op {
	case "Register":
		if Repo == nil { // Repo not initialized
			Repo = make(map[string]ior.IOR)
		}
		_p0 := request.Args[0].(string)
		_p1 := request.Args[1].(map[string]interface{})
		fmt.Printf("NamingServer:: %v\n", _p1)
		_p11 := _p1["Host"].(string)
		_p12 := _p1["Port"].(string)
		_p13 := _p1["Id"].(int64)
		_p14 := _p1["Proxy"].(string)

		iorTemp := ior.IOR{Host:_p11,Port:_p12,Id:int(_p13),Proxy:_p14}

		rTemp := Namingserver{}.Register(_p0,iorTemp)

		r := make([]interface{},2,2)
		r[0] = request.Op
		r[1] = rTemp
		*msg = messages.SAMessage{Payload:r}
	}
}

// TODO - REMOVE FROM HERE

type Naming struct{}

var Repo = map[string]ior.IOR{}

func (Namingserver) Lookup(s string) ior.IOR {
	return Repo[s]
}

func (Namingserver) List() []string{
	keys := make([]string, 0, len(Repo))
	for k := range Repo {
		keys = append(keys, k)
	}
	return keys
}

func (Namingserver) Register(serviceName string, ior ior.IOR) bool{
	if _, ok := Repo[serviceName]; ok {
		return false
	} else {
		Repo[serviceName] = ior
		return true
	}
}

