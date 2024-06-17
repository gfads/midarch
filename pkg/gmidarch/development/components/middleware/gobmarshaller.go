package middleware

import (
	"bytes"
	"encoding/gob"

	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages/miop"
	"github.com/gfads/midarch/pkg/shared"
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
		if req.Params[0] == nil {
			msg.Payload = messages.FunctionalReply{Rep: nil}
		} else {
			temp := req.Params[0].([]byte)
			r, err := g.Unmarshall(temp)
			if err != nil {
				msg.Payload = messages.FunctionalReply{Rep: nil}
			} else {
				msg.Payload = messages.FunctionalReply{Rep: r}
			}
		}
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

func (Gobmarshaller) Unmarshall(m []byte) (miop.MiopPacket, error) {
	r := new(miop.MiopPacket)
	var b bytes.Buffer
	b.Write(m)

	dec := gob.NewDecoder(&b)
	err := dec.Decode(r)
	if err != nil {
		// shared.ErrorHandler(shared.GetFunction(), err.Error())
		return miop.MiopPacket{}, err
	}

	// TODO improve by avoiding the loop
	// for i := 0; i < len(r.Bd.ReqBody.Body); i++ {
	// 	temp := r.Bd.ReqBody.Body[i]
	// 	if reflect.TypeOf(temp).String() == "float64" {
	// 		x := int(temp.(float64))
	// 		r.Bd.ReqBody.Body[i] = x
	// 	}
	// }

	return *r, nil
}
