package middleware

import (
	"bytes"
	"encoding/gob"
	"github.com/gfads/midarch/src/gmidarch/development/messages"
	"github.com/gfads/midarch/src/gmidarch/development/messages/miop"
	"github.com/gfads/midarch/src/shared"
)

// @Type: Gobmarshaller
// @Behaviour: Behaviour = InvP.e1 -> I_Process -> TerP.e1 -> Behaviour
type Gobmarshaller struct{}

func (g Gobmarshaller) I_Process(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	gob.Register(messages.FunctionalRequest{})
	gob.Register(messages.AOR{}) // TODO - perhaps put in init()

	req := msg.Payload.(messages.FunctionalRequest)
	op := req.Op

	switch op {
	case "marshall":
		r := g.Marshall(req.Params[0])
		msg.Payload = messages.FunctionalReply{Rep: r}
	case "unmarshall":
		temp := req.Params[0].([]byte)
		r := g.Unmarshall(temp)
		msg.Payload = messages.FunctionalReply{Rep: r}
	default:
		shared.ErrorHandler(shared.GetFunction(), "Operation '"+op+"' not supported by Marshaller!")
	}
}

func (Gobmarshaller) Marshall(m interface{}) []byte {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(m)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	return b.Bytes()
}
func (Gobmarshaller) Unmarshall(m []byte) interface{} {
	r := new(miop.MiopPacket)
	var b bytes.Buffer
	b.Write(m)

	dec := gob.NewDecoder(&b)
	err := dec.Decode(r)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	return *r
}
