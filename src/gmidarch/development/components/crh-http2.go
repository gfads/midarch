package components

import (
	"crypto/tls"
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"net/http"
	"os"
	"shared"
)

type CRHHttp2 struct {
	Behaviour string
	Graph     graphs.ExecGraph
	Conns     map[string]*tls.Conn
}

func NewCRHHttp2() CRHHttp2 {

	r := new(CRHHttp2)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"
	r.Conns = make(map[string]*tls.Conn, shared.NUM_MAX_CONNECTIONS)

	return *r
}

func (CRHHttp2) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	elem.(CRHHttp2).I_Process(msg, info)
}

func (c CRHHttp2) I_Process(msg *messages.SAMessage, info [] *interface{}) {
	// check message
	payload := msg.Payload.([]interface{})
	//host := payload[0].(string)        // host
	//port := payload[1].(string)        // port
	request := payload[2].(*http.Request) // HttpsRequest

	//addr := host + ":" + port
	var err error
	//if _, ok := c.Conns[key]; !ok { // no connection open yet
		//servAddr := key // TODO
		//tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
		//tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
		//if err != nil {
		//	log.Fatalf("CRHHttp2:: %s", err)
		//}

		//c.Conns[addr], err = tls.Dial("tcp4", addr, getClientTLSConfig())
		//if err != nil {
		//	fmt.Printf("CRHHttp2:: %v\n", err)
		//	os.Exit(1)
		//}
	//}

	// connect to server
	//conn := c.Conns[addr]

	// send message's size
	//size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	//binary.LittleEndian.PutUint32(size, uint32(len(msgToServer)))
	//_, err = conn.Write(size)
	//if err != nil {
	//	fmt.Printf("CRHHttp2:: %v\n", err)
	//	os.Exit(1)
	//}

	//fmt.Printf("CRHHttp2:: %v \n\n",size)

	// send message
	//fmt.Printf("CRHHttp2:: Message to server:: %v %v >> %v << \n\n",msgToServer, len(msgToServer), binary.LittleEndian.Uint32(size))
	//request := messages.HttpRequest{}
	//request.Unmarshal(msgToServer)

	client := &http.Client{}
	resp, err := client.Do(request)

	//resp, err := http.Get("https://"+addr+"/?"+request.QueryParameters)
	//_, err = conn.Write(msgToServer)
	if err != nil {
		fmt.Printf("CRHHttp2:: %v\n", err)
		os.Exit(1)
	}

	//fmt.Printf("CRHHttp2:: Message sent to Server [%v,%v] \n",conn.LocalAddr(),conn.RemoteAddr())

	// receive reply's size
	//_, err = conn.Read(size)
	//if err != nil {
	//	fmt.Printf("CRHHttp2:: %v\n", err)
	//	os.Exit(1)
	//}

	// receive reply
	//reader := bufio.NewReader(conn)
	//var message string
	//for {
	//	line, err := reader.ReadString('\n')
	//	if err == io.EOF {
	//		{
	//			fmt.Printf("SRHHttps:: Read\n")
	//			//os.Exit(0) // Todo: Comment this line. The server closes the connection after sending message. Maybe it's not necessary for performance evaluation
	//			break
	//		}
	//	} else if err != nil && err != io.EOF {
	//		fmt.Printf("SRHHttps:: %v\n", err)
	//		os.Exit(1)
	//	}
	//
	//	if strings.TrimSpace(line) == "" {
	//		break
	//	}
	//
	//	//fmt.Println("Request:", line, "END")
	//	message += line
	//}

	//response := messages.HttpResponse{}
	//response.Unmarshal(message)
	//cl, err := strconv.Atoi(response.Header.Fields["content-length"])
	//if err != nil {
	//	cl = 0
	//}
	//
	//if cl != 0 {
	//	body := make([]byte, cl, cl)
	//	reader.Read(body)
	//	response.Body = string(body)
	//}

	//fmt.Printf("CRHHttp2:: Message received from Server:: [%v,%v] \n",conn.LocalAddr(),conn.RemoteAddr())

	*msg = messages.SAMessage{Payload: resp}
}