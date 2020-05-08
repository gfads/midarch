package components

import (
	"encoding/binary"
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"io"
	"log"
	"net"
	"os"
	"shared"
)

type SRH struct {
	Behaviour string
	Graph     graphs.ExecGraph

	ConnsSRH 	[]net.Conn
	LnSRH 		net.Listener

	c1 					chan []byte
	c2 					chan []byte
	currentConnection 	int
	state 				int
}

//var ConnsSRH 	[]net.Conn
//var LnSRH 		net.Listener


//var c1 = make(chan []byte)
//var c2 = make(chan []byte)
//var currentConnection = -1
//var state = 0
var srh *SRH = &SRH{}

func NewSRH() SRH {
	//r SRH := null
	if srh.Behaviour == "" {
		fmt.Println("Vai instanciar o SRH pela primeira e única vez")
		srh = new(SRH)
		srh.Behaviour = "B = I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B"

		srh.c1 = make(chan []byte)
		srh.c2 = make(chan []byte)
		srh.currentConnection = -1
		srh.state = 0
	}
	fmt.Println("Passou no NewSRH")
	return *srh
}

func (s SRH) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'R' { // I_Receive
		srh.I_Receive(msg, info, elemInfo)
	} else { // "I_Send"
		srh.I_Send(msg, info, elemInfo)
	}
}

func (s SRH) I_Receive(msg *messages.SAMessage, info [] *interface{}, elemInfo [] *interface{}) { // TODO Host & Port
	tempPort := *elemInfo[0]
	port := tempPort.(string)
	host := "127.0.0.1" // TODO

	if s.LnSRH == nil { // listener was not created yet
		servAddr, err := net.ResolveTCPAddr("tcp", host+":"+port)
		if err != nil {
			log.Fatalf("SRH:: %v\n", err)
		}
		s.LnSRH, err = net.ListenTCP("tcp", servAddr)
		if err != nil {
			log.Fatalf("SRH:: %v\n", err)
		}
	}
	fmt.Println("State1:", s.state, " currentConnection1: ", s.currentConnection)
	switch s.state {
	case 0:
		go s.acceptAndRead(s.currentConnection, s.c1)
		s.state = 1
	case 1:
		go s.acceptAndRead(s.currentConnection, s.c1)
		s.state = 2
	case 2:
		go s.read(s.currentConnection, s.c2)
	}

	//go acceptAndRead(currentConnection, c1, done)
	//go read(currentConnection, c2, done)

	select {
	case msgTemp := <-s.c1:
		*msg = messages.SAMessage{Payload: msgTemp}
	case msgTemp := <-s.c2:
		*msg = messages.SAMessage{Payload: msgTemp}
	}

	s.currentConnection = s.nextConnection()
	fmt.Println("State1:", s.state, " currentConnection1: ", s.currentConnection)
}

func (s SRH) acceptAndRead(currentConnection int, c chan []byte) {

	// accept connections
	temp, err := s.LnSRH.Accept()
	if err != nil {
		fmt.Printf("SRH:: %v\n", err)
		os.Exit(1)
	}
	s.ConnsSRH = append(s.ConnsSRH, temp)
	s.currentConnection++

	// receive size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	tempConn := s.ConnsSRH[s.currentConnection]
	_, err = tempConn.Read(size)
	if err == io.EOF {
		{
			fmt.Printf("SRH:: Accept and Read\n")
			os.Exit(0)
		}
	} else if err != nil && err != io.EOF {
		fmt.Printf("SRH:: %v\n", err)
		os.Exit(1)
	}

	// receive message
	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	_, err = tempConn.Read(msgTemp)
	if err != nil {
		fmt.Printf("SRH:: %v\n", err)
		os.Exit(1)
	}
	c <- msgTemp
}

func (s SRH) read(currentConnection int, c chan []byte) {

	// receive size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	tempConn := s.ConnsSRH[s.currentConnection]

	_, err := tempConn.Read(size)
	if err == io.EOF {
		fmt.Printf("SRH:: Read\n")
		os.Exit(0)
	} else if err != nil && err != io.EOF {
		fmt.Printf("SRH:: %v\n", err)
		os.Exit(1)
	}

	// receive message
	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	_, err = tempConn.Read(msgTemp)
	if err != nil {
		fmt.Printf("SRH:: %v\n", err)
		os.Exit(1)
	}

	c <- msgTemp

	return
}

func (s SRH) I_Send(msg *messages.SAMessage, info [] *interface{}, elemInfo []*interface{}) {
	msgTemp := msg.Payload.([]interface{})[0].([]byte)
	fmt.Println(&s, "currentConnection:", s.currentConnection)
	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgTemp)))
	fmt.Println("ConnsSRH:", len(s.ConnsSRH))
	_, err := s.ConnsSRH[s.currentConnection].Write(size)
	fmt.Println("ConnsSRH 1.1")
	if err != nil {
		fmt.Printf("SRH:: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("ConnsSRH 1.2")
	// send message
	_, err = s.ConnsSRH[s.currentConnection].Write(msgTemp)
	fmt.Println("ConnsSRH 1.3")
	if err != nil {
		fmt.Printf("SRH:: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("ConnsSRH 1.4")
}

func (s SRH) nextConnection() int {
	fmt.Println(s, "entrou em nextConnection")
	r := -1

	if s.currentConnection == -1 {
		r = 0
	} else {
		if (s.currentConnection + 1) == len(s.ConnsSRH) {
			r = 0
		} else {
			r = s.currentConnection + 1
		}
	}
	return r
}