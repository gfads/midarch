package components

import (
	"encoding/binary"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"log"
	"net"
	"shared"
	"strconv"
)

type SRH struct {
	Behaviour string
	Graph     graphs.ExecGraph
	Host      string
	Port      int
	Conns     map[string]net.Conn
	Lns       map[string]net.Listener
}

func NewSRH() SRH {

	// create a new instance of Server
	r := new(SRH)

	// configure the new instance
	r.Host = "localhost"           // TODO
	r.Port = shared.FIBONACCI_PORT // TODO
	r.Behaviour = "B = I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B"
	r.Conns = make(map[string]net.Conn)
	r.Lns = make(map[string]net.Listener)

	return *r
}

func (s SRH) I_Receive(msg *messages.SAMessage, info [] *interface{}) { // TODO

	host := s.Host // TODO
	port := strconv.Itoa(s.Port) // TODO

	// create listener if necessary
	key := host + port
	err := *new(error)
	if _, ok := s.Lns[key]; !ok { // listener was not created yet
		servAddr, err := net.ResolveTCPAddr("tcp", host+":"+port)
		if err != nil {
			log.Fatalf("SRH:: %v\n", err)
		}
		s.Lns[key], err = net.ListenTCP("tcp", servAddr)
		if err != nil {
			log.Fatalf("SRH:: %v\n", err)
		}

		// accept connections
		s.Conns[key], err = s.Lns[key].Accept()
		if err != nil {
			log.Fatalf("SRH:: %s", err)
		}
	}

	// configure conn to be used
	conn := s.Conns[key]

	// receive size & message
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE)
	_, err = conn.Read(size)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	_, err = conn.Read(msgTemp)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	*msg = messages.SAMessage{Payload: msgTemp} // TODO
}

func (s SRH) I_Send(msg *messages.SAMessage, info [] *interface{}) {
	msgTemp := msg.Payload.([]interface{})[0].([]byte)

	// configure conn to be used
	key := s.Host + strconv.Itoa(s.Port)
	conn := s.Conns[key]

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgTemp)))
	_, err := conn.Write(size)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// send message
	_, err = conn.Write(msgTemp)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}
}
