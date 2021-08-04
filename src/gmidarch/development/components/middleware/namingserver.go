package middleware

import (
	"gmidarch/development/messages"
	"shared"
)

//@Type: Namingserver
//@Behaviour: Behaviour = InvP.e1 -> I_Process -> TerP.e1 -> Behaviour
type Namingserver struct{}

func (s Namingserver) I_Process(id string, msg *messages.SAMessage, info *interface{}) {
	request := msg.Payload.(*messages.FunctionalRequest)

	switch request.Op {
	case "Register":
		reply := Namingserver{}.Register(request.Params[0].(string), request.Params[1].(messages.AOR))
		msg.Payload = messages.FunctionalReply{Rep: reply}
	case "Lookup":
		reply, _ := Namingserver{}.Lookup(request.Params[0].(string))
		msg.Payload = messages.FunctionalReply{Rep: reply}
	case "List":
		reply := Namingserver{}.List()
		msg.Payload = messages.FunctionalReply{Rep:reply}
	default:
		shared.ErrorHandler(shared.GetFunction(), "Operation '"+request.Op+"' not implemented by Calculator")
	}
}

// TODO - REMOVE FROM HERE

type Naming struct{}

var Repo = map[string]messages.AOR{}

func (Namingserver) Register(serviceName string, aor messages.AOR) bool {
	r := false

	if _, ok := Repo[serviceName]; !ok {
		Repo[serviceName] = aor
		r = true
	}
	return r
}

func (Namingserver) Lookup(s string) (interface{}, bool) { // TODO dcruzb: dont need to be an interface, change to messages.AOR
	aor, ok := Repo[s]
	return aor, ok
}

func (Namingserver) List() []interface{} {
	keys := make([]interface{}, 0, len(Repo))
	for k := range Repo {
		keys = append(keys, k)
	}
	return keys
}
