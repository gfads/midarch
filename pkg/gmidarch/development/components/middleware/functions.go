package middleware

import (
	"net"
	"strings"

	"github.com/gfads/midarch/pkg/gmidarch/development/generic"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages/miop"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
	"github.com/quic-go/quic-go"
)

func VerifyProtocolAdaptation(msgFromServer []byte, sizeOfMsgSize []byte, protocol generic.Protocol) (err error) {
	if changeProtocol, miopPacket := isAdapt(msgFromServer); changeProtocol {
		lib.PrintlnDebug("Adapting, miopPacket.Bd.ReqBody.Body:", miopPacket.Bd.ReqBody.Body)
		shared.AdaptId = miopPacket.Bd.ReqBody.Body[1].(int)
		adaptToProtocol := miopPacket.Bd.ReqBody.Body[0].(string)
		confirmProtocolAdaptation(shared.AdaptId, adaptToProtocol, sizeOfMsgSize, protocol)
		prepareToAdaptTo(adaptToProtocol)
	}

	return nil
}

func VerifyAdaptation(msgFromServer []byte, sizeOfMsgSize []byte, conn net.Conn, send func(sizeOfMsgSize []byte, msgToServer []byte, conn net.Conn) error) (err error) {
	if changeProtocol, miopPacket := isAdapt(msgFromServer); changeProtocol {
		lib.PrintlnDebug("Adapting, miopPacket.Bd.ReqBody.Body:", miopPacket.Bd.ReqBody.Body)
		shared.AdaptId = miopPacket.Bd.ReqBody.Body[1].(int)
		protocol := miopPacket.Bd.ReqBody.Body[0].(string)
		confirmAdaptation(shared.AdaptId, protocol, sizeOfMsgSize, conn, send)
		prepareToAdaptTo(protocol)
	}

	return nil
}

func VerifyAdaptationQUIC(msgFromServer []byte, sizeOfMsgSize []byte, stream quic.Stream, send func(sizeOfMsgSize []byte, msgToServer []byte, stream quic.Stream) error) (err error) {
	if changeProtocol, miopPacket := isAdapt(msgFromServer); changeProtocol {
		lib.PrintlnDebug("Adapting, miopPacket.Bd.ReqBody.Body:", miopPacket.Bd.ReqBody.Body)
		shared.AdaptId = miopPacket.Bd.ReqBody.Body[1].(int)
		protocol := miopPacket.Bd.ReqBody.Body[0].(string)
		confirmAdaptationQUIC(shared.AdaptId, protocol, sizeOfMsgSize, stream, send)
		prepareToAdaptTo(protocol)
	}

	return nil
}

func isAdapt(msgFromServer []byte) (bool, miop.MiopPacket) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	miop := Jsonmarshaller{}.Unmarshall(msgFromServer)
	return miop.Bd.ReqHeader.Operation == "ChangeProtocol", miop
}

func confirmProtocolAdaptation(adaptId int, adaptToProtocol string, sizeOfMsgSize []byte, protocol generic.Protocol) (err error) {
	miopPacket := miop.CreateReqPacket("ChangeProtocol", []interface{}{adaptToProtocol, adaptId, "Ok"}, adaptId)
	msgPayload := Jsonmarshaller{}.Marshall(miopPacket)
	return protocol.Send(sizeOfMsgSize, msgPayload)
}

func confirmAdaptation(adaptId int, protocol string, sizeOfMsgSize []byte, conn net.Conn, send func(sizeOfMsgSize []byte, msgToServer []byte, conn net.Conn) error) (err error) {
	miopPacket := miop.CreateReqPacket("ChangeProtocol", []interface{}{protocol, adaptId, "Ok"}, adaptId)
	msgPayload := Jsonmarshaller{}.Marshall(miopPacket)
	return send(sizeOfMsgSize, msgPayload, conn)
}

func confirmAdaptationQUIC(adaptId int, protocol string, sizeOfMsgSize []byte, stream quic.Stream, send func(sizeOfMsgSize []byte, msgToServer []byte, stream quic.Stream) error) (err error) {
	miopPacket := miop.CreateReqPacket("ChangeProtocol", []interface{}{protocol, adaptId, "Ok"}, adaptId)
	msgPayload := Jsonmarshaller{}.Marshall(miopPacket)
	return send(sizeOfMsgSize, msgPayload, stream)
}

func prepareToAdaptTo(protocol string) (err error) {
	lib.PrintlnInfo("Adapting =>", strings.ToUpper(protocol))
	if protocol == "udp" {
		shared.ListOfComponentsToAdaptTo = append(shared.ListOfComponentsToAdaptTo, "crhudp")
	} else if protocol == "tcp" {
		shared.ListOfComponentsToAdaptTo = append(shared.ListOfComponentsToAdaptTo, "crhtcp")
	} else if protocol == "tls" {
		shared.ListOfComponentsToAdaptTo = append(shared.ListOfComponentsToAdaptTo, "crhtls")
	} else if protocol == "quic" {
		shared.ListOfComponentsToAdaptTo = append(shared.ListOfComponentsToAdaptTo, "crhquic")
	} else if protocol == "rpc" {
		shared.ListOfComponentsToAdaptTo = append(shared.ListOfComponentsToAdaptTo, "crhrpc")
	}

	return nil
}
