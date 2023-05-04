package middleware

import (
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages/miop"
	"github.com/gfads/midarch/pkg/shared"
)

// @Type: Requestor
// @Behaviour: Behaviour = InvP.e1 -> I_Beforemarshalling -> InvR.e2 -> TerR.e2 -> I_Beforesend -> InvR.e3 -> TerR.e3 -> I_Beforeunmarshalling -> InvR.e2 -> TerR.e2 -> I_Beforeproxy -> TerP.e1 -> Behaviour
type Requestor struct{}

func (Requestor) I_Beforemarshalling(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {

	//fmt.Println(shared.GetFunction(),msg)

	// Create invocation (to CRH/SRH) and Configure Info
	invocation := msg.Payload.(messages.Invocation)
	*info = messages.RequestorInfo{Inv: invocation}

	//Create request
	request := invocation.Functionalrequest
	reqPacket := miop.CreateReqPacket(request.Op, request.Params, shared.AdaptId)
	tempParams := []interface{}{reqPacket}

	// Configure message to marshaller
	msg.Payload = messages.FunctionalRequest{Op: "marshall", Params: tempParams}
}

func (Requestor) I_Beforesend(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {

	// Response from Marhsaller (only bytes)
	msgBytes := msg.Payload.(messages.FunctionalReply).Rep.([]byte)
	aux1 := *info
	aux2 := aux1.(messages.RequestorInfo)
	aux2.MarshalledMessage = msgBytes
	*info = aux2

	// Take the Response from marshaller
	msg.Payload = *info
}

func (Requestor) I_Beforeunmarshalling(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	tempParams := []interface{}{msg.Payload}
	msg.Payload = messages.FunctionalRequest{Op: "unmarshall", Params: tempParams}
}

func (Requestor) I_Beforeproxy(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//fmt.Println(shared.GetFunction(),msg)

	if msg.Payload.(messages.FunctionalReply).Rep == nil {
		msg.Payload = messages.FunctionalReply{Rep: nil}
	} else {
		temp1 := msg.Payload.(messages.FunctionalReply).Rep.(miop.MiopPacket)
		msg.Payload = messages.FunctionalReply{Rep: temp1.Bd.RepBody.OperationResult}
	}
}
