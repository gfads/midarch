package protocols

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
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

type HTTPSClient struct {
	connection net.Conn
	Ip         string
	adaptId    int
	msgChan    chan []byte
	replyChan  chan []byte
}

func (cl *HTTPSClient) Address() string {
	return cl.Ip
}

func (cl *HTTPSClient) AdaptId() int {
	return cl.adaptId
}

func (cl *HTTPSClient) SetAdaptId(adaptId int) {
	cl.adaptId = adaptId
}

func (cl *HTTPSClient) Connection() interface{} {
	return cl.connection
}

func (cl *HTTPSClient) CloseConnection() {
	cl.Ip = ""
	if cl.connection != nil {
		err := cl.connection.Close()
		if err != nil {
			lib.PrintlnError(err)
		}
	}
}

func (cl *HTTPSClient) ReadString() (message string) {
	//TODO implement me
	panic("implement me")
}

func (cl *HTTPSClient) WriteString(message string) {
	//TODO implement me
	panic("implement me")
}

func (cl *HTTPSClient) Read(b []byte) (err error) {
	//TODO implement me
	panic("implement me")
}

func (cl *HTTPSClient) Receive() (msg []byte, err error) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHHTTPS Version Not adapted")
	msg = <-cl.msgChan
	lib.PrintlnInfo("HTTPSClient.Receive: msg", msg)
	lib.PrintlnInfo("HTTPSClient.Receive: msg as string", string(msg))
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

func (cl *HTTPSClient) Send(msg []byte) error {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHHTTPS Version Not adapted")
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

type HTTPS struct {
	// Server attributes
	ip                 string
	port               string
	httpsServer        *http.Server
	listener           net.Listener
	initialConnections int
	clients            []*generic.Client
	started            bool

	// Client attributes
	httpsClient *http.Client // serverConnection net.Conn
	msgChan     chan []byte
}

// Server methods

func (st *HTTPS) StartServer(ip, port string, initialConnections int) {
	st.ip = ip
	st.port = port
	st.initialConnections = initialConnections

	lib.PrintlnInfo("HTTPS clients len", len(st.clients))
	if len(st.clients) < 1 { //st.initialConnections { TODO dcruzb : verify if there is the need to more than one client on HTTPS
		client := &HTTPSClient{}
		client.msgChan = make(chan []byte)
		client.replyChan = make(chan []byte)
		// *clientsPtr = append(clients, &client)
		st.AddClient(client, -1)
		lib.PrintlnInfo("HTTPS client created")
	}

	var client *HTTPSClient = (*st.clients[0]).(*HTTPSClient)
	request := new(HTTPSRequest)
	request.msgChan = client.msgChan
	request.replyChan = client.replyChan

	// Publish the receivers methods
	// handler := http.Handler{ServeHTTP: RequestTest}
	st.httpsServer = &http.Server{Handler: request}
	// st.httpsServer.Handler = request
	// st.httpsServer.Serve()

	// err := st.httpsServer.Handler(RequestTest)
	// if err != nil {
	// 	lib.PrintlnError(shared.GetFunction(), "Error while registering methods. Details:", err.Error())
	// }
	// Register a HTTPS handler
	// http.HandleHTTP()

	ln, err := tls.Listen("tcp4", st.ip+":"+st.port, lib.GetServerTLSConfig())
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), "Error while starting HTTPS server. Details: ", err)
	}

	st.listener = ln
}

func (st *HTTPS) StopServer() {
	err := st.listener.Close()
	if err != nil {
		lib.PrintlnError("Error while stoping server. Details:", err)
	}
}

func (st *HTTPS) AvailableConnectionFromPool() (available bool, idx int) {
	return !st.started, 0
}

func (st *HTTPS) WaitForConnection(cliIdx int) (cl *generic.Client) {
	st.started = true
	go func() {
		err := st.httpsServer.Serve(st.listener)
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
			shared.ErrorHandler(shared.GetFunction(), "Error while waiting for connection: "+err.Error())
		}
	}()

	lib.PrintlnInfo("HTTPS wait -> clients len", len(st.clients))
	if len(st.clients) > cliIdx {
		// (*st.clients[cliIdx]).(*HTTPSClient).connection = conn
		// (*st.clients[cliIdx]).(*HTTPSClient).Ip = conn.RemoteAddr().String()
		lib.PrintlnInfo("HTTPS wait -> client returned")
		return st.clients[cliIdx]
	} else {
		return nil
	}
}

func (st *HTTPS) GetClients() (client []*generic.Client) {
	return st.clients
}

func (st *HTTPS) GetClient(idx int) (client generic.Client) {
	return *st.clients[idx]
}

func (st *HTTPS) AddClient(client generic.Client, idx int) {
	if idx < 0 {
		st.clients = append(st.clients, &client)
	} else if idx < st.initialConnections {
		st.clients[idx] = &client
	}
}

func (st *HTTPS) GetClientFromAddr(addr string) (client generic.Client) {
	for _, client := range st.clients {
		if (*client).Address() == addr {
			return *client
		}
	}

	log.Println("IP without client from the ip:", addr)
	return nil
}

func (st *HTTPS) ResetClients() {
	for _, client := range st.clients {
		if client != nil {
			(*client).CloseConnection()
		}
	}
	st.clients = st.clients[:0]
}

// Client Methods

func (st *HTTPS) ConnectToServer(ip, port string) {
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
		st.httpsClient = &http.Client{Timeout: 5 * time.Second}
		st.httpsClient.Transport = &http.Transport{TLSClientConfig: getClientTLSConfig(), DialTLS: func(network, addr string) (net.Conn, error) { return tls.Dial(network, addr, getClientTLSConfig()) }}
		// st.serverConnection, err = net.DialTCP("tcp", nil, tcpAddr)
		lib.PrintlnDebug("Dialed", st.httpsClient)
		// if err != nil {
		// 	lib.PrintlnError("Dial error", st.httpsClient, err)
		// 	time.Sleep(200 * time.Millisecond)
		// 	//shared.ErrorHandler(shared.GetFunction(), err.Error())
		// } else {
		break // TODO dcruzb: remove for since there is no possibility of error
		// }
	}
	lib.PrintlnDebug("Connected", st.httpsClient)
	// if addr != shared.NAMING_HOST+":"+shared.NAMING_PORT && shared.LocalAddr == "" {
	// 	//lib.PrintlnDebug("crhInfo.Conns[addr].LocalAddr().String()", crhInfo.Conns[addr].LocalAddr().String())
	// 	shared.LocalAddr = st.httpsClient.LocalAddr().String()
	// 	lib.PrintlnDebug("Got local addr", st.httpsClient)
	// }
}

func (st *HTTPS) CloseConnection() {
	st.httpsClient.CloseIdleConnections()
	// err := st.httpsClient.Close()
	// if err != nil {
	// 	lib.PrintlnError(err)
	// }
}

func (st *HTTPS) ReadString() (message string) {
	//TODO implement me
	panic("implement me")
}

func (st *HTTPS) WriteString(message string) {
	//TODO implement me
	panic("implement me")
}

func (st *HTTPS) Receive() ([]byte, error) {
	lib.PrintlnInfo("----------------------------------------->", shared.GetFunction(), "HTTPS.Receive")
	msgFromServer := <-st.msgChan
	// sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	// // receive reply's size
	// _, err := st.serverConnection.Read(sizeOfMsgSize)
	// if err != nil {
	// 	lib.PrintlnError(shared.GetFunction(), err)
	// 	//shared.ErrorHandler(shared.GetFunction(), err.Error())
	// 	return nil, err
	// }
	// lib.PrintlnInfo("----------------------------------------->", shared.GetFunction(), "HTTPS read size")
	// // receive reply
	// msgFromServer := make([]byte, binary.LittleEndian.Uint32(sizeOfMsgSize), shared.NUM_MAX_MESSAGE_BYTES)
	// _, err = st.serverConnection.Read(msgFromServer)
	// if err != nil {
	// 	lib.PrintlnError(shared.GetFunction(), err)
	// 	//shared.ErrorHandler(shared.GetFunction(), err.Error())
	// 	return nil, err
	// }
	// lib.PrintlnInfo("----------------------------------------->", shared.GetFunction(), "HTTPS read message")
	return msgFromServer, nil
}

// func (st *HTTPS) Call(serviceMethod string, args any, reply any) (err error) {
// 	c := make(chan error, 1)
// 	go func() { c <- st.httpsClient.Call(serviceMethod, args, reply) }()
// 	select {
// 	case err = <-c:
// 		return err
// 	case <-time.After(2 * time.Second):
// 		return errors.New("HTTPS Call timeout")
// 	}
// }

func (st *HTTPS) Send(msgToServer []byte) error {
	lib.PrintlnInfo("CRHHTTPS Version Not adapted")
	//sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE) // TODO dcruzb: create attribute to avoid doing this everytime

	addr := st.ip + ":" + st.port

	// The message received from the server
	var msgFromServer []byte // := make([]byte, binary.LittleEndian.Uint32(sizeOfMsgSize), shared.NUM_MAX_MESSAGE_BYTES)
	response, err := st.httpsClient.Get("https://" + addr + "?param=" + url.PathEscape(string(msgToServer)))
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

type HTTPSRequest struct {
	msgChan   chan []byte
	replyChan chan []byte
	//ServeHTTP func(w http.ResponseWriter, r *http.Request)
}

func (rq HTTPSRequest) Request(w http.ResponseWriter, r *http.Request) { //request []byte, reply *[]byte) error {
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

func (rq HTTPSRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) { //request []byte, reply *[]byte) error {
	lib.PrintlnInfo("Received message. URI:", r.RequestURI)
	uriParameters := lib.GetURIParameters(r.RequestURI)
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

func getClientTLSConfig() *tls.Config {
	if shared.CA_PATH == "" {
		log.Fatal("CRHSsl:: Error:: Environment variable 'CA_PATH' not configured\n")
	}
	trustCert, err := ioutil.ReadFile(shared.CA_PATH)
	if err != nil {
		fmt.Println("Error loading trust certificate. ", err)
	}
	certs := x509.NewCertPool()
	if !certs.AppendCertsFromPEM(trustCert) {
		fmt.Println("Error installing trust certificate.")
	}

	// connect to server
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            certs,
		NextProtos:         []string{"h2"},
	}
	return tlsConfig
}
