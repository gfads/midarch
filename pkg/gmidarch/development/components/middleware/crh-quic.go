package middleware

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages/miop"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
	"github.com/quic-go/quic-go"
)

// @Type: CRHQUIC
// @Behaviour: Behaviour = InvP.e1 -> I_Process -> TerP.e1 -> Behaviour
type CRHQuic struct {
	//Conns map[string]quic.Connection
}

// var Stream quic.Stream

func (c CRHQuic) I_Process(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	infoTemp := *info
	crhInfo := infoTemp.(messages.CRHInfo)

	// check message
	payload := msg.Payload.(messages.RequestorInfo).MarshalledMessage
	h := msg.Payload.(messages.RequestorInfo).Inv.Endpoint.Host
	p := msg.Payload.(messages.RequestorInfo).Inv.Endpoint.Port
	msgToServer := payload

	host := ""
	port := ""

	if h == "" || p == "" {
		host = crhInfo.EndPoint.Host
		port = crhInfo.EndPoint.Port
	} else {
		host = h
		port = p
	}

	addr := host + ":" + port
	var err error
	if _, ok := crhInfo.QuicConns[addr]; !ok { // no connection open yet
		//tcpAddr, err := net.ResolveTCPAddr("tcp", key)
		//if err != nil {
		//	log.Fatalf("CRHQuic:: %s", err)
		//}

		crhInfo.QuicConns[addr], err = quic.DialAddr(addr, getClientTLSQuicConfig(), nil)
		if err != nil {
			fmt.Printf("CRHQuic:: %v\n", err)
			os.Exit(1)
		}

		crhInfo.QuicStreams[addr], err = crhInfo.QuicConns[addr].OpenStreamSync(context.Background())
		if err != nil {
			fmt.Printf("CRHQuic:: %v\n", err)
			os.Exit(1)
		}
	}

	// connect to server
	//conn := c.QuicConns[addr]

	// send message's size
	stream := crhInfo.QuicStreams[addr]
	sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	c.send(sizeOfMsgSize, msgToServer, stream)
	if err != nil {
		lib.PrintlnError("Error trying to send message:", err.Error())
		*msg = messages.SAMessage{Payload: nil} // TODO dcruzb: adjust message
		crhInfo.QuicStreams[addr].Close()
		crhInfo.QuicConns[addr] = nil
		crhInfo.QuicStreams[addr] = nil
		delete(crhInfo.QuicConns, addr)
		delete(crhInfo.QuicStreams, addr)
		return
	}
	//fmt.Printf("CRHQuic:: Message sent to Server [%v,%v] \n",conn.LocalAddr(),conn.RemoteAddr())

	msgFromServer, err := c.read(stream, sizeOfMsgSize)
	if err != nil {
		lib.PrintlnError("Error trying to read message:", err.Error())
		*msg = messages.SAMessage{Payload: nil} // TODO dcruzb: adjust message
		crhInfo.QuicStreams[addr].Close()
		crhInfo.QuicConns[addr] = nil
		crhInfo.QuicStreams[addr] = nil
		delete(crhInfo.QuicConns, addr)
		delete(crhInfo.QuicStreams, addr)
		return
	}
	//fmt.Printf("CRHQuic:: Message received from Server:: [%v,%v] \n",conn.LocalAddr(),conn.RemoteAddr())

	if changeProtocol, miopPacket := c.isAdapt(msgFromServer); changeProtocol {
		lib.PrintlnDebug("Adapting, miopPacket.Bd.ReqBody.Body:", miopPacket.Bd.ReqBody.Body)

		shared.AdaptId = miopPacket.Bd.ReqBody.Body[1].(int)

		miopPacket := miop.CreateReqPacket("ChangeProtocol", []interface{}{miopPacket.Bd.ReqBody.Body[0], shared.AdaptId, "Ok"}, shared.AdaptId) // idx is the Connection ID
		msgPayload := Jsonmarshaller{}.Marshall(miopPacket)
		c.send(sizeOfMsgSize, msgPayload, stream)

		if miopPacket.Bd.ReqBody.Body[0] == "udp" {
			lib.PrintlnInfo("Adapting => UDP")
			//evolutive.GeneratePlugin("crhudp_v1", "crhudp", "crhudp_v1")
			shared.ListOfComponentsToAdaptTo = append(shared.ListOfComponentsToAdaptTo, "crhudp")
		} else if miopPacket.Bd.ReqBody.Body[0] == "tcp" {
			lib.PrintlnInfo("Adapting => TCP")
			//evolutive.GeneratePlugin("crhtcp_v1", "crhtcp", "crhtcp_v1")
			shared.ListOfComponentsToAdaptTo = append(shared.ListOfComponentsToAdaptTo, "crhtcp")
		} else if miopPacket.Bd.ReqBody.Body[0] == "tls" {
			lib.PrintlnInfo("Adapting => TLS")
			//evolutive.GeneratePlugin("crhtcp_v1", "crhtcp", "crhtcp_v1")
			shared.ListOfComponentsToAdaptTo = append(shared.ListOfComponentsToAdaptTo, "crhtls")
		} else if miopPacket.Bd.ReqBody.Body[0] == "quic" {
			lib.PrintlnInfo("Adapting => QUIC")
			//evolutive.GeneratePlugin("crhtcp_v1", "crhtcp", "crhtcp_v1")
			shared.ListOfComponentsToAdaptTo = append(shared.ListOfComponentsToAdaptTo, "crhquic")
		} else {
			msgFromServer, _ = c.read(stream, sizeOfMsgSize)
			//fmt.Println("=================> ############### ============> ########### TCP: Leu o read")
		}
	}

	*msg = messages.SAMessage{Payload: msgFromServer}
}

func (c CRHQuic) send(sizeOfMsgSize []byte, msgToServer []byte, stream quic.Stream) error {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHQuic Version Not adapted")
	binary.LittleEndian.PutUint32(sizeOfMsgSize, uint32(len(msgToServer)))
	_, err := stream.Write(sizeOfMsgSize)
	if err != nil {
		// fmt.Printf("CRHQuic:: %v\n", err)
		return err
	}

	//fmt.Printf("CRHQuic:: %v \n\n",size)
	// send message
	//fmt.Printf("CRHQuic:: Message to server:: %v %v >> %v << \n\n",msgToServer, len(msgToServer), binary.LittleEndian.Uint32(size))
	_, err = stream.Write(msgToServer)
	if err != nil {
		// fmt.Printf("CRHQuic:: %v\n", err)
		return err
	}

	return nil
}

func (c CRHQuic) read(stream quic.Stream, size []byte) ([]byte, error) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHQuic Version Not adapted")
	// receive reply's size
	_, err := stream.Read(size)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err)
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return nil, err
	}

	// receive reply
	msgFromServer := make([]byte, binary.LittleEndian.Uint32(size), shared.NUM_MAX_MESSAGE_BYTES)
	_, err = stream.Read(msgFromServer)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err)
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return nil, err
	}
	return msgFromServer, nil
}

func (c CRHQuic) isAdapt(msgFromServer []byte) (bool, miop.MiopPacket) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHQuic Version Not adapted")
	miop := Jsonmarshaller{}.Unmarshall(msgFromServer)
	return miop.Bd.ReqHeader.Operation == "ChangeProtocol", miop
}

func getClientTLSQuicConfig() *tls.Config {
	if shared.CA_PATH == "" {
		log.Fatal("CRHQuic:: Error:: Environment variable 'CA_PATH' not configured\n")
	}

	trustCert, err := ioutil.ReadFile(shared.CA_PATH)
	if err != nil {
		fmt.Println("Error loading trust certificate. ", err)
	}
	certs := x509.NewCertPool()
	if !certs.AppendCertsFromPEM(trustCert) {
		fmt.Println("Error installing trust certificate.")
	}

	tlsConfig := &tls.Config{
		//InsecureSkipVerify: true,
		RootCAs:    certs,
		NextProtos: []string{"MidArchQuic"}, // TODO: Verify what NextProtos should be
	}
	return tlsConfig
}
