package middleware

import (
	"encoding/json"
	"gmidarch/development/messages"
	"gmidarch/development/messages/miop"
	"reflect"
	"shared"
)

//@Type: Jsonmarshaller
//@Behaviour: Behaviour = InvP.e1 -> I_Process -> TerP.e1 -> Behaviour
type Jsonmarshaller struct{}

func (j Jsonmarshaller) I_Process(id string, msg *messages.SAMessage, info *interface{}) {
	req := msg.Payload.(messages.FunctionalRequest)
	op := req.Op

	switch op {
	case "marshall":
		r := j.Marshall(req.Params[0].(miop.MiopPacket))
		msg.Payload = messages.FunctionalReply{Rep: r}
	case "unmarshall":
		temp := req.Params[0].([]byte)
		r := j.Unmarshall(temp)
		msg.Payload = messages.FunctionalReply{Rep: r}
	default:
		shared.ErrorHandler(shared.GetFunction(), "Marshaller:: Operation '"+op+"' not supported!")
	}
}

func (Jsonmarshaller) Marshall(m miop.MiopPacket) []byte {

	r, err := json.Marshal(m)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	return r
}

func (Jsonmarshaller) Unmarshall(m []byte) miop.MiopPacket {

	r := miop.MiopPacket{}
	err := json.Unmarshal(m, &r)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// TODO improve by avoiding the loop
	for i := 0; i < len(r.Bd.ReqBody.Body); i++ {
		temp := r.Bd.ReqBody.Body[i]
		if reflect.TypeOf(temp).String() == "float64" {
			x := int(temp.(float64))
			r.Bd.ReqBody.Body[i] = x
		}
	}
	return r
}
