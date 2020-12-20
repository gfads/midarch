package components

import (
	"bufio"
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type SRHHttp struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

var ConnsSRHHttp []net.Conn
var LnSRHHttp net.Listener

var c1Http = make(chan []byte)
var c2Http = make(chan []byte)
var currentConnectionHttp = -1
var stateHttp = 0

func NewSRHHttp() SRHHttp {

	r := new(SRHHttp)
	r.Behaviour = "B = I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B"

	return *r
}

func (e SRHHttp) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'R' { // I_Receive
		elem.(SRHHttp).I_Receive(msg, info, elemInfo)
	} else { // "I_Send"
		elem.(SRHHttp).I_Send(msg, info, elemInfo)
	}
}

func (e SRHHttp) I_Receive(msg *messages.SAMessage, info [] *interface{}, elemInfo [] *interface{}) { // TODO Host & Port
	tempPort := *elemInfo[0]
	port := tempPort.(string)
	host := "127.0.0.1" // TODO

	if LnSRHHttp == nil { // listener was not created yet
		servAddr, err := net.ResolveTCPAddr("tcp", host+":"+port)
		if err != nil {
			log.Fatalf("SRHHttp:: %v\n", err)
		}
		LnSRHHttp, err = net.ListenTCP("tcp", servAddr)
		if err != nil {
			log.Fatalf("SRHHttp:: %v\n", err)
		}
	}

	switch stateHttp {
	case 0:
		go acceptAndReadHttp(currentConnectionHttp, c1Http)
		stateHttp = 1
	case 1:
		go acceptAndReadHttp(currentConnectionHttp, c1Http)
		stateHttp = 2
	case 2:
		go acceptAndReadHttp(currentConnectionHttp, c1Http)
	}

	//go acceptAndReadHttp(currentConnectionHttp, c1Http, done)
	//go readHttp(currentConnectionHttp, c2Http, done)

	select {
	case msgTemp := <-c1Http:
		*msg = messages.SAMessage{Payload: msgTemp}
	case msgTemp := <-c2Http:
		*msg = messages.SAMessage{Payload: msgTemp}
	}

	currentConnectionHttp = nextConnectionHttp()
}

func acceptAndReadHttp(currentConnectionHttp int, c chan []byte) {

	// accept connections
	temp, err := LnSRHHttp.Accept()
	if err != nil {
		fmt.Printf("SRHHttp:: %v\n", err)
		os.Exit(1)
	}
	ConnsSRHHttp = append(ConnsSRHHttp, temp)
	currentConnectionHttp++

	// read request
	reader := bufio.NewReader(ConnsSRHHttp[currentConnectionHttp])
	var message string
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			{
				fmt.Printf("SRHHttp:: Accept and Read\n")
				os.Exit(0) // Todo: In a Http server EOF means that the client hit esc while waiting for result, it's not a problem! The server must stay up! What to do in this case?
			}
		} else if err != nil && err != io.EOF {
			fmt.Printf("SRHHttp:: %v\n", err)
			os.Exit(1)
		}

		if strings.TrimSpace(line) == "" { // Todo: supposing a request without body, have to change latter to support a body in requests
			break
		}

		//fmt.Println("Request:", line, "END")
		message += line
	}

	//fmt.Println("Size:",size)
	//fmt.Println("Size:",fmt.Sprint(binary.LittleEndian.Uint32(size)))
	//
	//// receive message
	//msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	//_, err = tempConn.Read(msgTemp)
	//if err != nil {
	//	fmt.Printf("SRHHttp:: %v\n", err)
	//	os.Exit(1)
	//}

	c <- []byte (message)
}

func readHttp(currentConnectionHttp int, c chan []byte) {

	// receive size
	//size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	//tempConn := ConnsSRHHttp[currentConnectionHttp]

	//_, err := tempConn.Read(size)
	//if err == io.EOF {
	//	fmt.Printf("SRHHttp:: read\n")
	//	os.Exit(0)
	//} else if err != nil && err != io.EOF {
	//	fmt.Printf("SRHHttp:: %v\n", err)
	//	os.Exit(1)
	//}
	//
	//// receive message
	//msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	//_, err = tempConn.Read(msgTemp)
	//if err != nil {
	//	fmt.Printf("SRHHttp:: %v\n", err)
	//	os.Exit(1)
	//}

	// read request
	reader := bufio.NewReader(ConnsSRHHttp[currentConnectionHttp])
	var message string
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			{
				fmt.Printf("SRHHttp:: Read\n")
				os.Exit(0) // Todo: In a Http server EOF means that the client hit esc while waiting for result, it's not a problem! The server must stay up! What to do in this case?
			}
		} else if err != nil && err != io.EOF {
			fmt.Printf("SRHHttp:: %v\n", err)
			os.Exit(1)
		}

		if strings.TrimSpace(line) == "" { // Todo: supposing a request without body, have to change latter to support a body in requests
			break
		}

		//fmt.Println("Request:", line, "END")
		message += line
	}

	c <- []byte (message)

	return
}

func (e SRHHttp) I_Send(msg *messages.SAMessage, info [] *interface{}, elemInfo []*interface{}) {
	msgTemp := msg.Payload.([]byte)


	// send message's size
	//size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	//binary.LittleEndian.PutUint32(size, uint32(len(msgTemp)))
	//_, err := ConnsSRHHttp[currentConnectionHttp].Write(size)
	//if err != nil {
	//	fmt.Printf("SRHHttp:: %v\n", err)
	//	os.Exit(1)
	//}

	// send message
	_, err := ConnsSRHHttp[currentConnectionHttp].Write(msgTemp)
	if err != nil {
		fmt.Printf("SRHHttp:: %v\n", err)
		os.Exit(1)
	}

	ConnsSRHHttp[currentConnectionHttp].Close()
}

func nextConnectionHttp() int {
	r := -1

	if currentConnectionHttp == -1 {
		r = 0
	} else {
		if (currentConnectionHttp + 1) == len(ConnsSRHHttp) {
			r = 0
		} else {
			r = currentConnectionHttp + 1
		}
	}
	return r
}
