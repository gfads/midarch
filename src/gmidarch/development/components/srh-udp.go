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

type SRHUdp struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

var connSRHUdp 	*net.UDPConn
var addr		*net.UDPAddr

var c1Udp = make(chan []byte)
var c2Udp = make(chan []byte)
//var currentConnectionUdp = -1
var stateUdp = 0

func NewSRHUdp() SRHUdp {

	r := new(SRHUdp)
	r.Behaviour = "B = I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B"

	return *r
}

func (s SRHUdp) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'R' { // I_Receive
		elem.(SRHUdp).I_Receive(msg, info, elemInfo)
	} else { // "I_Send"
		elem.(SRHUdp).I_Send(msg, info, elemInfo)
	}
}

func (s SRHUdp) I_Receive(msg *messages.SAMessage, info [] *interface{}, elemInfo [] *interface{}) { // TODO Host & Port
	tempPort := *elemInfo[0]
	port := tempPort.(string)
	host := "0.0.0.0" // TODO

	if connSRHUdp == nil { // connection was not created yet
		servAddr, err := net.ResolveUDPAddr("udp4", host+":"+port)
		if err != nil {
			log.Fatalf("SRHUdp:: %v\n", err)
		}
		connSRHUdp, err = net.ListenUDP("udp4", servAddr)

		if err != nil {
			log.Fatalf("SRHUdp:: %v\n", err)
		}
		//currentConnectionUdp++
	}

	// TODO: do these states make sense in UDP?
	switch stateUdp {
	case 0:
		go acceptAndReadUDP(0, c1Udp)
		stateUdp = 1
	case 1:
		go readUDP(0, c1Udp)
		stateUdp = 2
	case 2:
		go readUDP(0, c1Udp)
	}

	//go acceptAndRead(currentConnectionUdp, c1Udp, done)
	//go read(currentConnectionUdp, c2Udp, done)

	select {
	case msgTemp := <-c1Udp:
		*msg = messages.SAMessage{Payload: msgTemp}
	case msgTemp := <-c2Udp:
		*msg = messages.SAMessage{Payload: msgTemp}
	}

	//currentConnectionUdp = nextConnectionUDP()
}

// TODO: this function makes sense? Since UDP don't have connections its duplicated with readUDP
// TODO: remove currentConnection parameter since there is now connection array
func acceptAndReadUDP(currentConnection int, c chan []byte) {

	// accept connections
	//temp, err := LnSRH.Accept()
	//if err != nil {
	//	fmt.Printf("SRHUdp:: %v\n", err)
	//	os.Exit(1)
	//}
	//ConnsSRHUdp = append(ConnsSRHUdp, temp)
	//currentConnectionUdp++

	// receive size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	tempConn := connSRHUdp
	_, tmpAddr, err := tempConn.ReadFromUDP(size)
	addr = tmpAddr
	if err == io.EOF {
		{
			fmt.Printf("SRHUdp:: Accept and Read\n")
			os.Exit(0)
		}
	} else if err != nil && err != io.EOF {
		fmt.Printf("SRHUdp:: %v\n", err)
		os.Exit(1)
	}

	// receive message
	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	_, _, err = tempConn.ReadFromUDP(msgTemp)
	//addr = tmpAddr
	if err != nil {
		fmt.Printf("SRHUdp:: %v\n", err)
		os.Exit(1)
	}
	c <- msgTemp
}

func readUDP(currentConnection int, c chan []byte) {

	// receive size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	tempConn := connSRHUdp

	_, tmpAddr, err := tempConn.ReadFromUDP(size)
	addr = tmpAddr
	if err == io.EOF {
		fmt.Printf("SRHUdp:: Read\n")
		os.Exit(0)
	} else if err != nil && err != io.EOF {
		fmt.Printf("SRHUdp:: %v\n", err)
		os.Exit(1)
	}

	// receive message
	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	_, _, err = tempConn.ReadFromUDP(msgTemp)
	//addr = tmpAddr
	if err != nil {
		fmt.Printf("SRHUdp:: %v\n", err)
		os.Exit(1)
	}

	c <- msgTemp

	return
}

func (s SRHUdp) I_Send(msg *messages.SAMessage, info [] *interface{}, elemInfo []*interface{}) {
	msgTemp := msg.Payload.([]interface{})[0].([]byte)

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgTemp)))
	_, err := connSRHUdp.WriteToUDP(size, addr)
	if err != nil {
		fmt.Printf("SRHUdp:: %v\n", err)
		os.Exit(1)
	}

	// send message
	_, err = connSRHUdp.WriteToUDP(msgTemp, addr)
	if err != nil {
		fmt.Printf("SRHUdp:: %v\n", err)
		os.Exit(1)
	}
}

// TODO: what will identify the other peer if there is no connections?
func nextConnectionUDP() int {
	/*r := -1

	if currentConnectionUdp == -1 {
		r = 0
	} else {
		if (currentConnectionUdp + 1) == len(connSRHUdp) {
			r = 0
		} else {
			r = currentConnectionUdp + 1
		}
	}
	return r*/
	return 1
}
