package components

/*
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
}

var ConnsSRH map[string]net.Conn
var LnsSRH map[string]net.Listener
//var c1 = make(chan []byte)
//var c2 = make(chan []byte)

func NewSRH() SRH {

	r := new(SRH)

	r.Behaviour = "B = I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B"
	LnsSRH = make(map[string]net.Listener)
	ConnsSRH = make(map[string]net.Conn)

	return *r
}

func (e SRH) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'R' { // I_Receive
		elem.(SRH).I_Receive(msg, info, elemInfo)
	} else { // "I_Send"
		elem.(SRH).I_Send(msg, info, elemInfo)
	}
}

func (e SRH) I_Receive(msg *messages.SAMessage, info [] *interface{}, elemInfo [] *interface{}) { // TODO Host & Port
	tempPort := *elemInfo[0]
	port := tempPort.(string)
	host := "127.0.0.1" // TODO

	key := host + ":" + port
	if LnsSRH[key] == nil { // listener was not created yet
		servAddr, err := net.ResolveTCPAddr("tcp", key)
		if err != nil {
			fmt.Printf("SRH:: %s \n", err)
			os.Exit(0)
		}
		lnsTemp, err := net.ListenTCP("tcp", servAddr)
		if err != nil {
			fmt.Printf("SRH:: %s \n", err)
			os.Exit(0)
		}
		LnsSRH[key] = lnsTemp
	}

	// it allows to read/accept simultaneously
	var c1 = make(chan []byte)
	var c2 = make(chan []byte)

	go acceptAndRead(key, c1)
	if ConnsSRH[key] != nil {
		go read(key, c2)
	}

	select {
	case msgTemp := <-c1:
		*msg = messages.SAMessage{Payload: msgTemp}
	case msgTemp := <-c2:
		*msg = messages.SAMessage{Payload: msgTemp}
	}
}

func (e SRH) I_Send(msg *messages.SAMessage, info [] *interface{}, elemInfo []*interface{}) {
	msgTemp := msg.Payload.([]interface{})[0].([]byte)

	tempPort := *elemInfo[0]
	port := tempPort.(string)
	host := "127.0.0.1"
	key := host + ":" + port

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgTemp)))
	_, err := ConnsSRH[key].Write(size)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// send message
	_, err = ConnsSRH[key].Write(msgTemp)
	if err != nil {
		fmt.Printf("SRH:: %s \n", err)
		os.Exit(0)
	}
}

func acceptAndRead(key string, c chan []byte) {

	// accept connections
	var err error
	ConnsSRH[key], err = LnsSRH[key].Accept()
	if err != nil {
		fmt.Printf("SRH:: %s \n", err)
		os.Exit(0)
	}

	// receive size & message
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	tempConn := ConnsSRH[key]
	_, err = tempConn.Read(size)
	if err == io.EOF {
		{
			fmt.Printf("SRH:: Accept and Read\n")
			os.Exit(0)
		}
	} else if err != nil && err != io.EOF {
		fmt.Printf("SRH:: %s \n", err)
		os.Exit(0)
	}

	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	_, err = tempConn.Read(msgTemp)
	if err != nil {
		fmt.Printf("SRH:: %s \n", err)
		os.Exit(0)
	}

	c <- msgTemp
}

func read(key string, c chan []byte) {

	// receive size & message
	var err error
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	tempConn := ConnsSRH[key]
	_, err = tempConn.Read(size)
	if err == io.EOF {
		fmt.Printf("SRH:: Read\n")
		os.Exit(0)
	} else if err != nil && err != io.EOF {
		log.Fatalf("SRH::: %s", err)
	}

	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	_, err = tempConn.Read(msgTemp)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	c <- msgTemp
}
*/