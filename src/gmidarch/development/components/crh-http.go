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
	"shared"
	"strconv"
	"strings"
)

type CRHHttp struct {
	Behaviour string
	Graph     graphs.ExecGraph
	Conns     map[string]net.Conn
}

func NewCRHHttp() CRHHttp {

	r := new(CRHHttp)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"
	r.Conns = make(map[string]net.Conn, shared.NUM_MAX_CONNECTIONS)

	return *r
}

func (CRHHttp) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	elem.(CRHHttp).I_Process(msg, info)
}

func (c CRHHttp) I_Process(msg *messages.SAMessage, info [] *interface{}) {
	// check message
	payload := msg.Payload.([]interface{})
	host := payload[0].(string)        // host
	port := payload[1].(string)        // port
	msgToServer := payload[2].([]byte) // HttpRequest

	key := host + ":" + port
	var err error
	//if _, ok := c.Conns[key]; !ok { // no connection open yet
		//servAddr := key // TODO
		//tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
		tcpAddr, err := net.ResolveTCPAddr("tcp", key)
		if err != nil {
			log.Fatalf("CRHHttp:: %s", err)
		}

		c.Conns[key], err = net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			fmt.Printf("CRHHttp:: %v\n", err)
			os.Exit(1)
		}
	//}

	// connect to server
	conn := c.Conns[key]

	// send message's size
	//size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	//binary.LittleEndian.PutUint32(size, uint32(len(msgToServer)))
	//_, err = conn.Write(size)
	//if err != nil {
	//	fmt.Printf("CRHHttp:: %v\n", err)
	//	os.Exit(1)
	//}

	//fmt.Printf("CRHHttp:: %v \n\n",size)

	// send message
	//fmt.Printf("CRHHttp:: Message to server:: %v %v >> %v << \n\n",msgToServer, len(msgToServer), binary.LittleEndian.Uint32(size))
	request := messages.HttpRequest{}
	request.Unmarshal(msgToServer)

	_, err = conn.Write(msgToServer)
	if err != nil {
		fmt.Printf("CRHHttp:: %v\n", err)
		os.Exit(1)
	}

	//fmt.Printf("CRHHttp:: Message sent to Server [%v,%v] \n",conn.LocalAddr(),conn.RemoteAddr())

	// receive reply's size
	//_, err = conn.Read(size)
	//if err != nil {
	//	fmt.Printf("CRHHttp:: %v\n", err)
	//	os.Exit(1)
	//}

	// receive reply
	reader := bufio.NewReader(conn)
	var message string
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			{
				fmt.Printf("SRHHttp:: Read\n")
				//os.Exit(0) // Todo: Comment this line. The server closes the connection after sending message. Maybe it's not necessary for performance evaluation
				break
			}
		} else if err != nil && err != io.EOF {
			fmt.Printf("SRHHttp:: %v\n", err)
			os.Exit(1)
		}

		if strings.TrimSpace(line) == "" { 
			break
		}

		//fmt.Println("Request:", line, "END")
		message += line
	}

	response := messages.HttpResponse{}
	response.Unmarshal(message)
	cl, err := strconv.Atoi(response.Header.Fields["content-length"])
	if err != nil {
		cl = 0
	}

	if cl != 0 {
		body := make([]byte, cl, cl)
		reader.Read(body)
		response.Body = string(body)
	}

	//fmt.Printf("CRHHttp:: Message received from Server:: [%v,%v] \n",conn.LocalAddr(),conn.RemoteAddr())

	*msg = messages.SAMessage{Payload: response}
}