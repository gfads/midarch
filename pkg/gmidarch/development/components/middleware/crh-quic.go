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
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
	"github.com/quic-go/quic-go"
)

// @Type: CRHQUIC
// @Behaviour: Behaviour = InvP.e1 -> I_Process -> TerP.e1 -> Behaviour
type CRHQUIC struct {
	//Conns map[string]quic.Connection
}

// var Stream quic.Stream

func (c CRHQUIC) I_Process(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
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
		//	log.Fatalf("CRHQUIC:: %s", err)
		//}

		crhInfo.QuicConns[addr], err = quic.DialAddr(addr, getClientTLSQuicConfig(), nil)
		if err != nil {
			fmt.Printf("CRHQUIC:: %v\n", err)
			os.Exit(1)
		}

		crhInfo.QuicStreams[addr], err = crhInfo.QuicConns[addr].OpenStreamSync(context.Background())
		if err != nil {
			fmt.Printf("CRHQUIC:: %v\n", err)
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
	//fmt.Printf("CRHQUIC:: Message sent to Server [%v,%v] \n",conn.LocalAddr(),conn.RemoteAddr())

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
	//fmt.Printf("CRHQUIC:: Message received from Server:: [%v,%v] \n",conn.LocalAddr(),conn.RemoteAddr())

	VerifyAdaptationQUIC(msgFromServer, sizeOfMsgSize, stream, c.send)

	*msg = messages.SAMessage{Payload: msgFromServer}
}

func (c CRHQUIC) send(sizeOfMsgSize []byte, msgToServer []byte, stream quic.Stream) error {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHQUIC Version Not adapted")
	binary.LittleEndian.PutUint32(sizeOfMsgSize, uint32(len(msgToServer)))
	_, err := stream.Write(sizeOfMsgSize)
	if err != nil {
		// fmt.Printf("CRHQUIC:: %v\n", err)
		return err
	}

	//fmt.Printf("CRHQUIC:: %v \n\n",size)
	// send message
	//fmt.Printf("CRHQUIC:: Message to server:: %v %v >> %v << \n\n",msgToServer, len(msgToServer), binary.LittleEndian.Uint32(size))
	_, err = stream.Write(msgToServer)
	if err != nil {
		// fmt.Printf("CRHQUIC:: %v\n", err)
		return err
	}

	return nil
}

func (c CRHQUIC) read(stream quic.Stream, size []byte) ([]byte, error) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHQUIC Version Not adapted")
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

func getClientTLSQuicConfig() *tls.Config {
	if shared.CA_PATH == "" {
		log.Fatal("CRHQUIC:: Error:: Environment variable 'CA_PATH' not configured\n")
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
