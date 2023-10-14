package protocols

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gfads/midarch/pkg/gmidarch/development/generic"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
)

type HTTPClient struct {
	connection net.Conn
	Ip         string
	adaptId    int
	msgChan    chan []byte
	replyChan  chan []byte
}

func (cl *HTTPClient) Address() string {
	return cl.Ip
}

func (cl *HTTPClient) AdaptId() int {
	return cl.adaptId
}

func (cl *HTTPClient) SetAdaptId(adaptId int) {
	cl.adaptId = adaptId
}

func (cl *HTTPClient) Connection() interface{} {
	return cl.connection
}

func (cl *HTTPClient) CloseConnection() {
	cl.Ip = ""
	if cl.connection != nil {
		err := cl.connection.Close()
		if err != nil {
			lib.PrintlnError(err)
		}
	}
}

func (cl *HTTPClient) ReadString() (message string) {
	//TODO implement me
	panic("implement me")
}

func (cl *HTTPClient) WriteString(message string) {
	//TODO implement me
	panic("implement me")
}

func (cl *HTTPClient) Read(b []byte) (err error) {
	//TODO implement me
	panic("implement me")
}

func (cl *HTTPClient) Receive() (msg []byte, err error) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHHTTP Version Not adapted")
	msg = <-cl.msgChan
	lib.PrintlnInfo("HTTPClient.Receive: msg", msg)
	lib.PrintlnInfo("HTTPClient.Receive: msg as string", string(msg))
	// receive reply's size
	// size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	// cl.Read(size)
	// if err != nil {
	// 	lib.PrintlnError(shared.GetFunction(), err)
	// 	//shared.ErrorHandler(shared.GetFunction(), err.Error())
	// 	return nil, err
	// }
	// // receive reply
	// msg = make([]byte, binary.LittleEndian.Uint32(size), shared.NUM_MAX_MESSAGE_BYTES)
	// err = cl.Read(msg)
	// if err != nil {
	// 	lib.PrintlnError(shared.GetFunction(), err)
	// 	//shared.ErrorHandler(shared.GetFunction(), err.Error())
	// 	return nil, err
	// }
	return msg, nil
}

func (cl *HTTPClient) Send(msg []byte) error {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHHTTP Version Not adapted")
	go func() {
		cl.replyChan <- msg
	}()

	// sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE) // TODO dcruzb: create attribute to avoid doing this everytime
	// binary.LittleEndian.PutUint32(sizeOfMsgSize, uint32(len(msg)))
	// _, err := cl.connection.Write(sizeOfMsgSize)
	// if err != nil {
	// 	//shared.ErrorHandler(shared.GetFunction(), err.Error())
	// 	return err
	// }

	// // send message
	// _, err = cl.connection.Write(msg)
	// if err != nil {
	// 	//shared.ErrorHandler(shared.GetFunction(), err.Error())
	// 	return err
	// }
	return nil
}

type HTTP struct {
	// Server attributes
	ip                 string
	port               string
	httpServer         *http.Server
	listener           net.Listener
	initialConnections int
	clients            []*generic.Client
	started            bool

	// Client attributes
	httpClient *http.Client // serverConnection net.Conn
	msgChan    chan []byte
}

// Server methods

func (st *HTTP) StartServer(ip, port string, initialConnections int) {
	st.ip = ip
	st.port = port
	st.initialConnections = initialConnections

	lib.PrintlnInfo("HTTP clients len", len(st.clients))
	if len(st.clients) < 1 { //st.initialConnections { TODO dcruzb : verify if there is the need to more than one client on HTTP
		client := &HTTPClient{}
		client.msgChan = make(chan []byte)
		client.replyChan = make(chan []byte)
		// *clientsPtr = append(clients, &client)
		st.AddClient(client, -1)
		lib.PrintlnInfo("HTTP client created")
	}

	var client *HTTPClient = (*st.clients[0]).(*HTTPClient)
	request := new(HTTPRequest)
	request.msgChan = client.msgChan
	request.replyChan = client.replyChan

	// Publish the receivers methods
	// handler := http.Handler{ServeHTTP: RequestTest}
	st.httpServer = &http.Server{Handler: request}
	// st.httpServer.Handler = request
	// st.httpServer.Serve()

	// err := st.httpServer.Handler(RequestTest)
	// if err != nil {
	// 	lib.PrintlnError(shared.GetFunction(), "Error while registering methods. Details:", err.Error())
	// }
	// Register a HTTP handler
	// http.HandleHTTP()

	ln, err := net.Listen("tcp", st.ip+":"+st.port)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), "Error while starting HTTP server. Details: ", err)
	}

	st.listener = ln
}

func (st *HTTP) StopServer() {
	err := st.listener.Close()
	if err != nil {
		lib.PrintlnError("Error while stoping server. Details:", err)
	}
}

func (st *HTTP) AvailableConnectionFromPool() (available bool, idx int) {
	return !st.started, 0
}

func (st *HTTP) WaitForConnection(cliIdx int) (cl *generic.Client) {
	st.started = true
	go func() {
		err := st.httpServer.Serve(st.listener)
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
			shared.ErrorHandler(shared.GetFunction(), "Error while waiting for connection: "+err.Error())
		}
	}()

	lib.PrintlnInfo("HTTP wait -> clients len", len(st.clients))
	if len(st.clients) > cliIdx {
		// (*st.clients[cliIdx]).(*HTTPClient).connection = conn
		// (*st.clients[cliIdx]).(*HTTPClient).Ip = conn.RemoteAddr().String()
		lib.PrintlnInfo("HTTP wait -> client returned")
		return st.clients[cliIdx]
	} else {
		return nil
	}
}

func (st *HTTP) GetClients() (client []*generic.Client) {
	return st.clients
}

func (st *HTTP) GetClient(idx int) (client generic.Client) {
	return *st.clients[idx]
}

func (st *HTTP) AddClient(client generic.Client, idx int) {
	if idx < 0 {
		st.clients = append(st.clients, &client)
	} else if idx < st.initialConnections {
		st.clients[idx] = &client
	}
}

func (st *HTTP) GetClientFromAddr(addr string) (client generic.Client) {
	for _, client := range st.clients {
		if (*client).Address() == addr {
			return *client
		}
	}

	log.Println("IP without client from the ip:", addr)
	return nil
}

func (st *HTTP) ResetClients() {
	for _, client := range st.clients {
		if client != nil {
			(*client).CloseConnection()
		}
	}
	st.clients = st.clients[:0]
}

// Client Methods

func (st *HTTP) ConnectToServer(ip, port string) {
	lib.PrintlnInfo("**********************************************")
	if st.msgChan == nil {
		st.msgChan = make(chan []byte)
	}
	st.ip = ip
	st.port = port
	// addr := st.ip + ":" + st.port
	// tcpAddr, err := net.ResolveTCPAddr("tcp", addr)

	// if err != nil {
	// 	shared.ErrorHandler(shared.GetFunction(), err.Error())
	// }
	// lib.PrintlnDebug("Resolved addr", tcpAddr)
	//localTcpAddr := c.getLocalTcpAddr()

	for {
		// Create an HTTP client with a timeout
		st.httpClient = &http.Client{Timeout: 5 * time.Second}
		// st.serverConnection, err = net.DialTCP("tcp", nil, tcpAddr)
		lib.PrintlnDebug("Dialed", st.httpClient)
		// if err != nil {
		// 	lib.PrintlnError("Dial error", st.httpClient, err)
		// 	time.Sleep(200 * time.Millisecond)
		// 	//shared.ErrorHandler(shared.GetFunction(), err.Error())
		// } else {
		break // TODO dcruzb: remove for since there is no possibility of error
		// }
	}
	lib.PrintlnDebug("Connected", st.httpClient)
	// if addr != shared.NAMING_HOST+":"+shared.NAMING_PORT && shared.LocalAddr == "" {
	// 	//lib.PrintlnDebug("crhInfo.Conns[addr].LocalAddr().String()", crhInfo.Conns[addr].LocalAddr().String())
	// 	shared.LocalAddr = st.httpClient.LocalAddr().String()
	// 	lib.PrintlnDebug("Got local addr", st.httpClient)
	// }
}

func (st *HTTP) CloseConnection() {
	st.httpClient.CloseIdleConnections()
	// err := st.httpClient.Close()
	// if err != nil {
	// 	lib.PrintlnError(err)
	// }
}

func (st *HTTP) ReadString() (message string) {
	//TODO implement me
	panic("implement me")
}

func (st *HTTP) WriteString(message string) {
	//TODO implement me
	panic("implement me")
}

func (st *HTTP) Receive() ([]byte, error) {
	lib.PrintlnInfo("----------------------------------------->", shared.GetFunction(), "HTTP.Receive")
	msgFromServer := <-st.msgChan
	// sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	// // receive reply's size
	// _, err := st.serverConnection.Read(sizeOfMsgSize)
	// if err != nil {
	// 	lib.PrintlnError(shared.GetFunction(), err)
	// 	//shared.ErrorHandler(shared.GetFunction(), err.Error())
	// 	return nil, err
	// }
	// lib.PrintlnInfo("----------------------------------------->", shared.GetFunction(), "HTTP read size")
	// // receive reply
	// msgFromServer := make([]byte, binary.LittleEndian.Uint32(sizeOfMsgSize), shared.NUM_MAX_MESSAGE_BYTES)
	// _, err = st.serverConnection.Read(msgFromServer)
	// if err != nil {
	// 	lib.PrintlnError(shared.GetFunction(), err)
	// 	//shared.ErrorHandler(shared.GetFunction(), err.Error())
	// 	return nil, err
	// }
	// lib.PrintlnInfo("----------------------------------------->", shared.GetFunction(), "HTTP read message")
	return msgFromServer, nil
}

// func (st *HTTP) Call(serviceMethod string, args any, reply any) (err error) {
// 	c := make(chan error, 1)
// 	go func() { c <- st.httpClient.Call(serviceMethod, args, reply) }()
// 	select {
// 	case err = <-c:
// 		return err
// 	case <-time.After(2 * time.Second):
// 		return errors.New("HTTP Call timeout")
// 	}
// }

func (st *HTTP) Send(msgToServer []byte) error {
	lib.PrintlnInfo("CRHHTTP Version Not adapted")
	//sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE) // TODO dcruzb: create attribute to avoid doing this everytime

	addr := st.ip + ":" + st.port

	// The message received from the server
	var msgFromServer []byte // := make([]byte, binary.LittleEndian.Uint32(sizeOfMsgSize), shared.NUM_MAX_MESSAGE_BYTES)
	response, err := st.httpClient.Get("http://" + addr + "?param=" + url.PathEscape(string(msgToServer)))
	if err != nil {
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return err
	}
	defer response.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return err
	}

	// // Convert the byte slice to a string
	// responseBody := string(bodyBytes)

	msgFromServer = bodyBytes

	lib.PrintlnInfo("Got message from server" + string(msgFromServer))
	go func() {
		st.msgChan <- msgFromServer
	}()
	lib.PrintlnInfo("Put message in msgChan")
	return nil

	// binary.LittleEndian.PutUint32(sizeOfMsgSize, uint32(len(msgToServer)))
	// _, err := st.serverConnection.Write(sizeOfMsgSize)
	// if err != nil {
	// 	//shared.ErrorHandler(shared.GetFunction(), err.Error())
	// 	return err
	// }

	// // send message
	// _, err = st.serverConnection.Write(msgToServer)
	// if err != nil {
	// 	//shared.ErrorHandler(shared.GetFunction(), err.Error())
	// 	return err
	// }
	// return nil
}

type HTTPRequest struct {
	msgChan   chan []byte
	replyChan chan []byte
	//ServeHTTP func(w http.ResponseWriter, r *http.Request)
}

func (rq HTTPRequest) Request(w http.ResponseWriter, r *http.Request) { //request []byte, reply *[]byte) error {
	lib.PrintlnInfo("Received message")
	uriParameters := getURIParameters(r.RequestURI)
	go func() {
		rq.msgChan <- []byte(uriParameters["param"].(string))
	}()
	lib.PrintlnInfo("Forwarded message")
	replyMsg := <-rq.replyChan
	lib.PrintlnInfo("Received reply")
	//*reply = w
	w.Write(replyMsg)

	//return nil
}

func (rq HTTPRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) { //request []byte, reply *[]byte) error {
	lib.PrintlnInfo("Received message. URI:", r.RequestURI)
	uriParameters := getURIParameters(r.RequestURI)
	lib.PrintlnInfo("Received message:", uriParameters["param"].(string))
	go func() {
		rq.msgChan <- []byte(uriParameters["param"].(string))
	}()
	lib.PrintlnInfo("Forwarded message")
	replyMsg := <-rq.replyChan
	lib.PrintlnInfo("Received reply")
	//*reply = w
	w.Write(replyMsg)

	//return nil
}

func getURIParameters(uri string) (parameters map[string]interface{}) {
	decodedPathParam, err := url.PathUnescape(uri[8:])
	if err != nil {
		lib.PrintlnInfo("Error decoding path parameter:", err)
		return
	}
	return map[string]interface{}{"param": decodedPathParam}

	// parameters = make(map[string]interface{})
	// paramRegex, err := regexp.Compile("([?][\\w|=|%|(|)|+|-|.|:]+)|([&][\\w|=|%|(|)|+|-|.|:]+)")
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	// params := paramRegex.FindAllString(uri, -1)
	// for _, param := range params {
	// 	param := param[1:]
	// 	keyValue := strings.Split(param, "=")
	// 	parameters[keyValue[0]] = keyValue[1]
	// }
	// return parameters
}
