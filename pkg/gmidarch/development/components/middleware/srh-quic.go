package middleware

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"io"
	"log"
	"os"
	"shared"

	"github.com/quic-go/quic-go"
)

type SRHQuic struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

var ConnsSRHQuic []quic.Session
var StreamsSRHQuic []quic.Stream
var LnSRHQuic quic.Listener

var c1Quic = make(chan []byte)
var c2Quic = make(chan []byte)
var currentConnectionQuic = -1
var stateQuic = 0

func NewSRHQuic() SRHQuic {

	r := new(SRHQuic)
	r.Behaviour = "B = I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B"

	return *r
}

func (e SRHQuic) Selector(elem interface{}, elemInfo []*interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'R' { // I_Receive
		elem.(SRHQuic).I_Receive(msg, info, elemInfo)
	} else { // "I_Send"
		elem.(SRHQuic).I_Send(msg, info, elemInfo)
	}
}

func (e SRHQuic) I_Receive(msg *messages.SAMessage, info []*interface{}, elemInfo []*interface{}) { // TODO Host & Port
	tempPort := *elemInfo[0]
	port := tempPort.(string)
	host := "0.0.0.0" //"127.0.0.1" // TODO

	if LnSRHQuic == nil { // listener was not created yet
		//servAddr, err := net.ResolveTCPAddr("tcp", host+":"+port)
		//if err != nil {
		//	log.Fatalf("SRH:: %v\n", err)
		//}
		quicConfig := quic.Config{KeepAlive: true}
		ln, err := quic.ListenAddr(host+":"+port, getServerTLSQuicConfig(), &quicConfig)
		if err != nil {
			log.Fatalf("SRHQuic:: %v\n", err)
		}
		LnSRHQuic = ln
	}

	switch stateQuic {
	case 0:
		go acceptAndReadQuic(currentConnectionQuic, c1Quic)
		stateQuic = 1
	case 1:
		go readQuic(currentConnectionQuic, c1Quic)
		stateQuic = 2
	case 2:
		go readQuic(currentConnectionQuic, c1Quic)
	}

	//go acceptAndRead(currentConnectionQuic, c1Quic, done)
	//go read(currentConnectionQuic, c2Quic, done)

	select {
	case msgTemp := <-c1Quic:
		*msg = messages.SAMessage{Payload: msgTemp}
	case msgTemp := <-c2Quic:
		*msg = messages.SAMessage{Payload: msgTemp}
	}

	currentConnectionQuic = nextConnectionQuic()
}

func acceptAndReadQuic(currentConnectionQuic int, c chan []byte) {

	// accept connections
	temp, err := LnSRHQuic.Accept(context.Background())
	if err != nil {
		fmt.Printf("SRHQuic:: %v\n", err)
		os.Exit(1)
	}
	ConnsSRHQuic = append(ConnsSRHQuic, temp) // Quic Session
	currentConnectionQuic++

	// receive size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	tempConn := ConnsSRHQuic[currentConnectionQuic]
	stream, err := tempConn.AcceptStream(context.Background())
	//stream, err := tempConn.OpenStreamSync(context.Background())
	StreamsSRHQuic = append(StreamsSRHQuic, stream)
	if err != nil {
		fmt.Printf("SRHQuic:: %v\n", err)
		os.Exit(1)
	}
	_, err = stream.Read(size)
	if err == io.EOF {
		{
			fmt.Printf("SRHQuic:: Accept and Read\n")
			os.Exit(0)
		}
	} else if err != nil && err != io.EOF {
		fmt.Printf("SRHQuic:: %v\n", err)
		os.Exit(1)
	}
	stream2 := StreamsSRHQuic[currentConnectionQuic]
	// receive message
	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	_, err = stream2.Read(msgTemp)
	if err != nil {
		fmt.Printf("SRHQuic:: %v\n", err)
		os.Exit(1)
	}
	c <- msgTemp
}

func readQuic(currentConnectionQuic int, c chan []byte) {
	// receive size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	stream := StreamsSRHQuic[currentConnectionQuic]

	_, err := stream.Read(size)
	if err == io.EOF {
		fmt.Printf("SRHQuic:: Read\n")
		os.Exit(0)
	} else if err != nil && err != io.EOF {
		fmt.Printf("SRHQuic:: %v\n", err)
		os.Exit(1)
	}

	// receive message
	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	_, err = stream.Read(msgTemp)
	if err != nil {
		fmt.Printf("SRHQuic:: %v\n", err)
		os.Exit(1)
	}

	c <- msgTemp

	return
}

func (e SRHQuic) I_Send(msg *messages.SAMessage, info []*interface{}, elemInfo []*interface{}) {
	msgTemp := msg.Payload.([]interface{})[0].([]byte)

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgTemp)))
	_, err := StreamsSRHQuic[currentConnectionQuic].Write(size)
	if err != nil {
		fmt.Printf("SRHQuic:: %v\n", err)
		os.Exit(1)
	}

	// send message
	_, err = StreamsSRHQuic[currentConnectionQuic].Write(msgTemp)
	if err != nil {
		fmt.Printf("SRHQuic:: %v\n", err)
		os.Exit(1)
	}
}

func nextConnectionQuic() int {
	r := -1

	if currentConnectionQuic == -1 {
		r = 0
	} else {
		if (currentConnectionQuic + 1) == len(ConnsSRHQuic) {
			r = 0
		} else {
			r = currentConnectionQuic + 1
		}
	}
	return r
}

func getServerTLSQuicConfig() *tls.Config {
	if shared.CRT_PATH == "" {
		log.Fatal("SRHSsl:: Error:: Environment variable 'CRT_PATH' not configured\n")
	}

	if shared.KEY_PATH == "" {
		log.Fatal("SRHSsl:: Error:: Environment variable 'KEY_PATH' not configured\n")
	}

	cert, err := tls.LoadX509KeyPair(shared.CRT_PATH, shared.KEY_PATH)
	if err != nil {
		log.Fatal("Error loading certificate. ", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"MidArchQuic"}, // TODO: Verify what NextProtos should be
	}
	return tlsConfig
}
