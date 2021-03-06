package components

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"io/ioutil"
	"log"
	"os"
	"shared"
)

type CRHQuic struct {
	Behaviour string
	Graph     graphs.ExecGraph
	Conns     map[string]quic.Session
}

var Stream    quic.Stream

func NewCRHQuic() CRHQuic {
	r := new(CRHQuic)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"
	r.Conns = make(map[string]quic.Session, shared.NUM_MAX_CONNECTIONS)

	return *r
}

func (CRHQuic) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	elem.(CRHQuic).I_Process(msg, info)
}

func (c CRHQuic) I_Process(msg *messages.SAMessage, info [] *interface{}) {
	// check message
	payload := msg.Payload.([]interface{})
	host := payload[0].(string)        // host
	port := payload[1].(string)        // port
	msgToServer := payload[2].([]byte)

	addr := host + ":" + port
	var err error
	if _, ok := c.Conns[addr]; !ok { // no connection open yet
		//tcpAddr, err := net.ResolveTCPAddr("tcp", key)
		//if err != nil {
		//	log.Fatalf("CRHQuic:: %s", err)
		//}

		c.Conns[addr], err = quic.DialAddr(addr, getClientTLSQuicConfig(), nil)
		if err != nil {
			fmt.Printf("CRHQuic:: %v\n", err)
			os.Exit(1)
		}

		Stream, err = c.Conns[addr].OpenStreamSync(context.Background())
		if err != nil {
			fmt.Printf("CRHQuic:: %v\n", err)
			os.Exit(1)
		}
	}

	// connect to server
	//conn := c.Conns[addr]

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgToServer)))
	_, err = Stream.Write(size)
	if err != nil {
		fmt.Printf("CRHQuic:: %v\n", err)
		os.Exit(1)
	}
	// TODO continue from here

	//fmt.Printf("CRHQuic:: %v \n\n",size)
		// send message
	//fmt.Printf("CRHQuic:: Message to server:: %v %v >> %v << \n\n",msgToServer, len(msgToServer), binary.LittleEndian.Uint32(size))
	_, err = Stream.Write(msgToServer)
	if err != nil {
		fmt.Printf("CRHQuic:: %v\n", err)
		os.Exit(1)
	}

	//fmt.Printf("CRHQuic:: Message sent to Server [%v,%v] \n",conn.LocalAddr(),conn.RemoteAddr())

	// receive reply's size
	_, err = Stream.Read(size)
	if err != nil {
		fmt.Printf("CRHQuic:: %v\n", err)
		os.Exit(1)
	}


	// receive reply
	msgFromServer := make([]byte, binary.LittleEndian.Uint32(size), shared.NUM_MAX_MESSAGE_BYTES)
	_, err = Stream.Read(msgFromServer)
	if err != nil {
		fmt.Printf("CRHQuic:: %v\n", err)
		os.Exit(1)
	}
	//fmt.Printf("CRHQuic:: Message received from Server:: [%v,%v] \n",conn.LocalAddr(),conn.RemoteAddr())

	*msg = messages.SAMessage{Payload: msgFromServer}
}

func getClientTLSQuicConfig() *tls.Config {
	if shared.CA_PATH == "" {
		log.Fatal("CRHSsl:: Error:: Environment variable 'CA_PATH' not configured\n")
	}

	trustCert, err := ioutil.ReadFile(shared.CA_PATH)
	if err != nil {
		fmt.Println("Error loading trust certificate. ",err)
	}
	certs := x509.NewCertPool()
	if !certs.AppendCertsFromPEM(trustCert) {
		fmt.Println("Error installing trust certificate.")
	}

	tlsConfig := &tls.Config{
		//InsecureSkipVerify: true,
		RootCAs: certs,
		NextProtos:         []string{"MidArchQuic"}, // TODO: Verify what NextProtos should be
	}
	return tlsConfig
}