package middleware

import (
	sendFileImpl "github.com/gfads/midarch/examples/sendfiledistributed/sendfileImpl"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages/miop"
	"github.com/gfads/midarch/pkg/shared"
)

// @Type: SendFileInvoker
// @Behaviour: Behaviour = InvP.e1 -> I_Beforeunmarshalling -> InvR.e2 -> TerR.e2 -> I_Beforeserver -> I_Beforemarshalling -> InvR.e2 -> TerR.e2 -> I_Beforesend -> TerP.e1 -> Behaviour
type SendFileInvoker struct{}

func (SendFileInvoker) I_Beforeunmarshalling(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	tempParams := []interface{}{msg.Payload}
	msg.Payload = messages.FunctionalRequest{Op: "unmarshall", Params: tempParams}
}

func (SendFileInvoker) I_Beforeserver(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	if msg.Payload.(messages.FunctionalReply).Rep == nil {
		shared.ErrorHandler(shared.GetFunction(), "Payload is nil")
	}
	miopPacket := msg.Payload.(messages.FunctionalReply).Rep.(miop.MiopPacket) // from marshaller

	req := messages.FunctionalRequest{Op: miopPacket.Bd.ReqHeader.Operation, Params: miopPacket.Bd.ReqBody.Body}

	switch req.Op {
	case "U":
		reply := sendFileImpl.SendFile{}.Save(req.Params[0].([]byte)) //.(string))
		msg.Payload = messages.FunctionalReply{Rep: reply}
	default:
		shared.ErrorHandler(shared.GetFunction(), "Operation '"+req.Op+"' not present in Invoker")
	}
}

func (SendFileInvoker) I_Beforemarshalling(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	reply := msg.Payload.(messages.FunctionalReply)
	repPacket := miop.CreateRepPacket(reply.Rep)
	tempParams := []interface{}{repPacket}
	msg.Payload = messages.FunctionalRequest{Op: "marshall", Params: tempParams}
}

func (SendFileInvoker) I_Beforesend(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	msg.Payload = msg.Payload.(messages.FunctionalReply).Rep
}
