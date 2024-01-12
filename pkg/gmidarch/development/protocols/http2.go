package protocols

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gfads/midarch/pkg/gmidarch/development/generic"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
	"golang.org/x/net/http2"
)

type HTTP2Client struct {
	connection net.Conn
	Ip         string
	adaptId    int
	msgChan    chan []byte
	replyChan  chan []byte
}

func (cl *HTTP2Client) Address() string {
	return cl.Ip
}

func (cl *HTTP2Client) AdaptId() int {
	return cl.adaptId
}

func (cl *HTTP2Client) SetAdaptId(adaptId int) {
	cl.adaptId = adaptId
}

func (cl *HTTP2Client) Connection() interface{} {
	return cl.connection
}

func (cl *HTTP2Client) CloseConnection() {
	cl.Ip = ""
	if cl.connection != nil {
		err := cl.connection.Close()
		if err != nil {
			lib.PrintlnError(err)
		}
	}
}

func (cl *HTTP2Client) ReadString() (message string) {
	//TODO implement me
	panic("implement me")
}

func (cl *HTTP2Client) WriteString(message string) {
	//TODO implement me
	panic("implement me")
}

func (cl *HTTP2Client) Read(b []byte) (err error) {
	//TODO implement me
	panic("implement me")
}

func (cl *HTTP2Client) Receive() (msg []byte, err error) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHHTTP2 Version Not adapted")
	msg = <-cl.msgChan
	// lib.PrintlnInfo("HTTP2Client.Receive: msg", msg)
	// lib.PrintlnInfo("HTTP2Client.Receive: msg as string", string(msg))
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

func (cl *HTTP2Client) Send(msg []byte) error {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHHTTP2 Version Not adapted")
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

type HTTP2 struct {
	// Server attributes
	ip                 string
	port               string
	http2Server        *http.Server
	listener           net.Listener
	initialConnections int
	clients            []*generic.Client
	started            bool

	// Client attributes
	http2Client *http.Client // serverConnection net.Conn
	msgChan     chan []byte
}

// Server methods

func (st *HTTP2) StartServer(ip, port string, initialConnections int) {
	st.ip = ip
	st.port = port
	st.initialConnections = initialConnections

	lib.PrintlnInfo("HTTP2 clients len", len(st.clients))
	if len(st.clients) < 1 { //st.initialConnections { TODO dcruzb : verify if there is the need to more than one client on HTTP2
		client := &HTTP2Client{}
		client.msgChan = make(chan []byte)
		client.replyChan = make(chan []byte)
		// *clientsPtr = append(clients, &client)
		st.AddClient(client, -1)
		lib.PrintlnInfo("HTTP2 client created")
	}

	var client *HTTP2Client = (*st.clients[0]).(*HTTP2Client)
	request := new(HTTP2Request)
	request.msgChan = client.msgChan
	request.replyChan = client.replyChan

	// Publish the receivers methods
	// handler := http.Handler{ServeHTTP: RequestTest}
	st.http2Server = &http.Server{Handler: request}
	// st.http2Server.Handler = request
	// st.http2Server.Serve()

	// err := st.http2Server.Handler(RequestTest)
	// if err != nil {
	// 	lib.PrintlnError(shared.GetFunction(), "Error while registering methods. Details:", err.Error())
	// }
	// Register a HTTP2 handler
	// http.HandleHTTP()

	ln, err := tls.Listen("tcp4", st.ip+":"+st.port, lib.GetServerTLSConfig("h2")) // TODO dcruzb: use https but with protos as h2?
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), "Error while starting HTTP2 server. Details: ", err)
	}

	st.listener = ln
}

func (st *HTTP2) StopServer() {
	err := st.listener.Close()
	if err != nil {
		lib.PrintlnError("Error while stoping server. Details:", err)
	}
}

func (st *HTTP2) AvailableConnectionFromPool() (available bool, idx int) {
	return !st.started, 0
}

func (st *HTTP2) WaitForConnection(cliIdx int) (cl *generic.Client) {
	st.started = true
	go func() {
		err := st.http2Server.Serve(st.listener)
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
			shared.ErrorHandler(shared.GetFunction(), "Error while waiting for connection: "+err.Error())
		}
	}()

	lib.PrintlnInfo("HTTP2 wait -> clients len", len(st.clients))
	if len(st.clients) > cliIdx {
		// (*st.clients[cliIdx]).(*HTTP2Client).connection = conn
		// (*st.clients[cliIdx]).(*HTTP2Client).Ip = conn.RemoteAddr().String()
		lib.PrintlnInfo("HTTP2 wait -> client returned")
		return st.clients[cliIdx]
	} else {
		return nil
	}
}

func (st *HTTP2) GetClients() (client []*generic.Client) {
	return st.clients
}

func (st *HTTP2) GetClient(idx int) (client generic.Client) {
	return *st.clients[idx]
}

func (st *HTTP2) AddClient(client generic.Client, idx int) {
	if idx < 0 {
		st.clients = append(st.clients, &client)
	} else if idx < st.initialConnections {
		st.clients[idx] = &client
	}
}

func (st *HTTP2) GetClientFromAddr(addr string) (client generic.Client) {
	for _, client := range st.clients {
		if (*client).Address() == addr {
			return *client
		}
	}

	log.Println("IP without client from the ip:", addr)
	return nil
}

func (st *HTTP2) ResetClients() {
	for _, client := range st.clients {
		if client != nil {
			(*client).CloseConnection()
		}
	}
	st.clients = st.clients[:0]
}

// Client Methods

func (st *HTTP2) ConnectToServer(ip, port string) {
	lib.PrintlnInfo("********************************************** HTTP2.ConnectToServer")
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

	// Create an HTTP client with a timeout
	http2Transport := &http2.Transport{TLSClientConfig: lib.GetClientTLSConfig("h2")}
	st.http2Client = &http.Client{Timeout: 5 * time.Second, Transport: http2Transport}
	lib.PrintlnDebug("Connected", st.http2Client)
	// if addr != shared.NAMING_HOST+":"+shared.NAMING_PORT && shared.LocalAddr == "" {
	// 	//lib.PrintlnDebug("crhInfo.Conns[addr].LocalAddr().String()", crhInfo.Conns[addr].LocalAddr().String())
	// 	shared.LocalAddr = st.http2Client.LocalAddr().String()
	// 	lib.PrintlnDebug("Got local addr", st.http2Client)
	// }
}

func (st *HTTP2) CloseConnection() {
	st.http2Client.CloseIdleConnections()
	// err := st.http2Client.Close()
	// if err != nil {
	// 	lib.PrintlnError(err)
	// }
}

func (st *HTTP2) ReadString() (message string) {
	//TODO implement me
	panic("implement me")
}

func (st *HTTP2) WriteString(message string) {
	//TODO implement me
	panic("implement me")
}

func (st *HTTP2) Receive() ([]byte, error) {
	lib.PrintlnInfo("----------------------------------------->", shared.GetFunction(), "HTTP2.Receive")
	msgFromServer := <-st.msgChan
	// sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	// // receive reply's size
	// _, err := st.serverConnection.Read(sizeOfMsgSize)
	// if err != nil {
	// 	lib.PrintlnError(shared.GetFunction(), err)
	// 	//shared.ErrorHandler(shared.GetFunction(), err.Error())
	// 	return nil, err
	// }
	// lib.PrintlnInfo("----------------------------------------->", shared.GetFunction(), "HTTP2 read size")
	// // receive reply
	// msgFromServer := make([]byte, binary.LittleEndian.Uint32(sizeOfMsgSize), shared.NUM_MAX_MESSAGE_BYTES)
	// _, err = st.serverConnection.Read(msgFromServer)
	// if err != nil {
	// 	lib.PrintlnError(shared.GetFunction(), err)
	// 	//shared.ErrorHandler(shared.GetFunction(), err.Error())
	// 	return nil, err
	// }
	// lib.PrintlnInfo("----------------------------------------->", shared.GetFunction(), "HTTP2 read message")
	return msgFromServer, nil
}

// func (st *HTTP2) Call(serviceMethod string, args any, reply any) (err error) {
// 	c := make(chan error, 1)
// 	go func() { c <- st.http2Client.Call(serviceMethod, args, reply) }()
// 	select {
// 	case err = <-c:
// 		return err
// 	case <-time.After(2 * time.Second):
// 		return errors.New("HTTP2 Call timeout")
// 	}
// }

func (st *HTTP2) Send(msgToServer []byte) error {
	lib.PrintlnInfo("CRHHTTP2 Version Not adapted")
	//sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE) // TODO dcruzb: create attribute to avoid doing this everytime
	lib.PrintlnInfo("************************************************************************ 1")
	addr := st.ip + ":" + st.port

	// The message received from the server
	var msgFromServer []byte // := make([]byte, binary.LittleEndian.Uint32(sizeOfMsgSize), shared.NUM_MAX_MESSAGE_BYTES)
	lib.PrintlnInfo("************************************************************************ 2")
	req, err := http.NewRequest("GET", "https://"+addr, bytes.NewReader(msgToServer))
	// req.Header.Set("Accept-Encoding", "gzip")
	response, err := st.http2Client.Do(req)
	lib.PrintlnInfo("************************************************************************ 3")
	lib.PrintlnInfo("response:", response)
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

type HTTP2Request struct {
	msgChan   chan []byte
	replyChan chan []byte
	//ServeHTTP func(w http.ResponseWriter, r *http.Request)
}

func (rq HTTP2Request) Request(w http.ResponseWriter, r *http.Request) { //request []byte, reply *[]byte) error {
	lib.PrintlnInfo("Received message")
	uriParameters := lib.GetURIParameters(r.RequestURI)
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

func (rq HTTP2Request) ServeHTTP(w http.ResponseWriter, r *http.Request) { //request []byte, reply *[]byte) error {
	// lib.PrintlnInfo("Received message. URI:", r.RequestURI)
	// uriParameters := lib.GetURIParameters(r.RequestURI)
	// lib.PrintlnInfo("Received message:", uriParameters["param"].(string))
	// go func() {
	// 	rq.msgChan <- []byte(uriParameters["param"].(string))
	// }()
	msg, _ := io.ReadAll(r.Body)
	r.Body.Close()
	go func() {
		rq.msgChan <- msg //[]byte(uriParameters["param"].(string))
	}()
	lib.PrintlnInfo("Forwarded message")
	replyMsg := <-rq.replyChan
	lib.PrintlnInfo("Received reply")
	//*reply = w
	w.Write(replyMsg)

	//return nil
}
