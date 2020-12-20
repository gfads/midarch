package components

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type SRHHttps struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

var ConnsSRHHttps []net.Conn
var LnSRHHttps net.Listener

var c1Https = make(chan []byte)
var c2Https = make(chan []byte)
var currentConnectionHttps = -1
var stateHttps = 0

func NewSRHHttps() SRHHttps {

	r := new(SRHHttps)
	r.Behaviour = "B = I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B"

	return *r
}

func (e SRHHttps) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'R' { // I_Receive
		elem.(SRHHttps).I_Receive(msg, info, elemInfo)
	} else { // "I_Send"
		elem.(SRHHttps).I_Send(msg, info, elemInfo)
	}
}

func (e SRHHttps) I_Receive(msg *messages.SAMessage, info [] *interface{}, elemInfo [] *interface{}) { // TODO Host & Port
	tempPort := *elemInfo[0]
	port := tempPort.(string)
	host := "localhost" // TODO

	if LnSRHHttps == nil { // listener was not created yet
		//servAddr, err := net.ResolveTCPAddr("tcp", host+":"+port)
		//if err != nil {
		//	log.Fatalf("SRHHttps:: %v\n", err)
		//}
		ln, err := tls.Listen("tcp4", host+":"+port, getServerTLSConfig())
		if err != nil {
			log.Fatalf("SRHHttps:: %v\n", err)
		}
		LnSRHHttps = ln
	}

	switch stateHttps {
	case 0:
		go acceptAndReadHttps(currentConnectionHttps, c1Https)
		stateHttps = 1
	case 1:
		go acceptAndReadHttps(currentConnectionHttps, c1Https)
		stateHttps = 2
	case 2:
		go acceptAndReadHttps(currentConnectionHttps, c1Https)
	}

	//go acceptAndReadHttps(currentConnectionHttps, c1Https, done)
	//go readHttps(currentConnectionHttps, c2Https, done)

	select {
	case msgTemp := <-c1Https:
		*msg = messages.SAMessage{Payload: msgTemp}
	case msgTemp := <-c2Https:
		*msg = messages.SAMessage{Payload: msgTemp}
	}

	currentConnectionHttps = nextConnectionHttps()
}

func acceptAndReadHttps(currentConnectionHttps int, c chan []byte) {

	// accept connections
	temp, err := LnSRHHttps.Accept()
	if err != nil {
		fmt.Printf("SRHHttps:: %v\n", err)
		os.Exit(1)
	}
	ConnsSRHHttps = append(ConnsSRHHttps, temp)
	currentConnectionHttps++

	// read request
	reader := bufio.NewReader(ConnsSRHHttps[currentConnectionHttps])
	var message string
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			{
				fmt.Printf("SRHHttps:: Accept and Read\n")
				os.Exit(0) // Todo: In a Https server EOF means that the client hit esc while waiting for result, it's not a problem! The server must stay up! What to do in this case?
			}
		} else if err != nil && err != io.EOF {
			fmt.Printf("SRHHttps:: %v\n", err)
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
	//	fmt.Printf("SRHHttps:: %v\n", err)
	//	os.Exit(1)
	//}

	c <- []byte (message)
}

func readHttps(currentConnectionHttps int, c chan []byte) {

	// receive size
	//size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	//tempConn := ConnsSRHHttps[currentConnectionHttps]

	//_, err := tempConn.Read(size)
	//if err == io.EOF {
	//	fmt.Printf("SRHHttps:: read\n")
	//	os.Exit(0)
	//} else if err != nil && err != io.EOF {
	//	fmt.Printf("SRHHttps:: %v\n", err)
	//	os.Exit(1)
	//}
	//
	//// receive message
	//msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	//_, err = tempConn.Read(msgTemp)
	//if err != nil {
	//	fmt.Printf("SRHHttps:: %v\n", err)
	//	os.Exit(1)
	//}

	// read request
	reader := bufio.NewReader(ConnsSRHHttps[currentConnectionHttps])
	var message string
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			{
				fmt.Printf("SRHHttps:: Read\n")
				os.Exit(0) // Todo: In a Https server EOF means that the client hit esc while waiting for result, it's not a problem! The server must stay up! What to do in this case?
			}
		} else if err != nil && err != io.EOF {
			fmt.Printf("SRHHttps:: %v\n", err)
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

func (e SRHHttps) I_Send(msg *messages.SAMessage, info [] *interface{}, elemInfo []*interface{}) {
	msgTemp := msg.Payload.([]byte)


	// send message's size
	//size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	//binary.LittleEndian.PutUint32(size, uint32(len(msgTemp)))
	//_, err := ConnsSRHHttps[currentConnectionHttps].Write(size)
	//if err != nil {
	//	fmt.Printf("SRHHttps:: %v\n", err)
	//	os.Exit(1)
	//}

	// send message
	_, err := ConnsSRHHttps[currentConnectionHttps].Write(msgTemp)
	if err != nil {
		fmt.Printf("SRHHttps:: %v\n", err)
		os.Exit(1)
	}

	ConnsSRHHttps[currentConnectionHttps].Close()
}

func nextConnectionHttps() int {
	r := -1

	if currentConnectionHttps == -1 {
		r = 0
	} else {
		if (currentConnectionHttps + 1) == len(ConnsSRHHttps) {
			r = 0
		} else {
			r = currentConnectionHttps + 1
		}
	}
	return r
}
