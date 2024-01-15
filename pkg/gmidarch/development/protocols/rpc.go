package protocols

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strings"
	"time"

	"github.com/gfads/midarch/pkg/gmidarch/development/generic"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
)

type RPCClient struct {
	connection net.Conn
	Ip         string
	adaptId    int
	msgChan    chan []byte
	replyChan  chan []byte
}

func (cl *RPCClient) Address() string {
	return cl.Ip
}

func (cl *RPCClient) AdaptId() int {
	return cl.adaptId
}

func (cl *RPCClient) SetAdaptId(adaptId int) {
	cl.adaptId = adaptId
}

func (cl *RPCClient) Connection() interface{} {
	return cl.connection
}

func (cl *RPCClient) CloseConnection() {
	cl.Ip = ""
	if cl.connection != nil {
		err := cl.connection.Close()
		if err != nil {
			lib.PrintlnError(err)
		}
	}
}

func (cl *RPCClient) ReadString() (message string) {
	//TODO implement me
	panic("implement me")
}

func (cl *RPCClient) WriteString(message string) {
	//TODO implement me
	panic("implement me")
}

func (cl *RPCClient) Read(b []byte) (err error) {
	//TODO implement me
	panic("implement me")
}

func (cl *RPCClient) Receive() (msg []byte, err error) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHRPC Version Not adapted")
	msg = <-cl.msgChan
	// lib.PrintlnInfo("RPCClient.Receive: msg", msg)
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

func (cl *RPCClient) Send(msg []byte) error {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHRPC Version Not adapted")
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

type RPC struct {
	// Server attributes
	ip                 string
	port               string
	rpcServer          *rpc.Server
	listener           net.Listener
	initialConnections int
	clients            []*generic.Client
	started            bool

	// Client attributes
	rpcClient *rpc.Client // serverConnection net.Conn
	msgChan   chan []byte
}

// Server methods

func (st *RPC) StartServer(ip, port string, initialConnections int) {
	st.ip = ip
	st.port = port
	st.initialConnections = initialConnections

	// lib.PrintlnInfo("RPC clients len", len(st.clients))
	if len(st.clients) < 1 { //st.initialConnections { TODO dcruzb : verify if there is the need to more than one client on RPC
		client := &RPCClient{}
		client.msgChan = make(chan []byte)
		client.replyChan = make(chan []byte)
		// *clientsPtr = append(clients, &client)
		st.AddClient(client, -1)
		lib.PrintlnInfo("RPC client created")
	}

	var client *RPCClient = (*st.clients[0]).(*RPCClient)
	request := new(RPCRequest)
	request.msgChan = client.msgChan
	request.replyChan = client.replyChan

	// Publish the receivers methods
	st.rpcServer = rpc.NewServer()

	err := st.rpcServer.Register(request)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), "Error while registering methods. Details:", err.Error())
	}
	// Register a HTTP handler
	// rpc.HandleHTTP()

	ln, err := net.Listen("tcp", st.ip+":"+st.port)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), "Error while starting RPC server. Details: ", err)
	}

	st.listener = ln
}

func (st *RPC) StopServer() {
	err := st.listener.Close()
	if err != nil {
		lib.PrintlnError("Error while stoping server. Details:", err)
	}
}

func (st *RPC) AvailableConnectionFromPool() (available bool, idx int) {
	return !st.started, 0
}

func (st *RPC) WaitForConnection(cliIdx int) (cl *generic.Client) {
	st.started = true
	go func() {
		err := http.Serve(st.listener, st.rpcServer)
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
			shared.ErrorHandler(shared.GetFunction(), "Error while waiting for connection: "+err.Error())
		}
	}()

	// lib.PrintlnInfo("RPC wait -> clients len", len(st.clients))
	if len(st.clients) > cliIdx {
		// (*st.clients[cliIdx]).(*RPCClient).connection = conn
		// (*st.clients[cliIdx]).(*RPCClient).Ip = conn.RemoteAddr().String()
		// lib.PrintlnInfo("RPC wait -> client returned")
		return st.clients[cliIdx]
	} else {
		return nil
	}
}

func (st *RPC) GetClients() (client []*generic.Client) {
	return st.clients
}

func (st *RPC) GetClient(idx int) (client generic.Client) {
	return *st.clients[idx]
}

func (st *RPC) AddClient(client generic.Client, idx int) {
	if idx < 0 {
		st.clients = append(st.clients, &client)
	} else if idx < st.initialConnections {
		st.clients[idx] = &client
	}
}

func (st *RPC) GetClientFromAddr(addr string) (client generic.Client) {
	for _, client := range st.clients {
		if (*client).Address() == addr {
			return *client
		}
	}

	log.Println("IP without client from the ip:", addr)
	return nil
}

func (st *RPC) ResetClients() {
	for _, client := range st.clients {
		if client != nil {
			(*client).CloseConnection()
		}
	}
	st.clients = st.clients[:0]
}

// Client Methods

func (st *RPC) ConnectToServer(ip, port string) {
	// lib.PrintlnInfo("**********************************************")
	if st.msgChan == nil {
		st.msgChan = make(chan []byte)
	}
	st.ip = ip
	st.port = port
	addr := st.ip + ":" + st.port
	// tcpAddr, err := net.ResolveTCPAddr("tcp", addr)

	// if err != nil {
	// 	shared.ErrorHandler(shared.GetFunction(), err.Error())
	// }
	// lib.PrintlnDebug("Resolved addr", tcpAddr)
	//localTcpAddr := c.getLocalTcpAddr()

	for {
		rpcClient, err := rpc.DialHTTP("tcp", addr)
		st.rpcClient = rpcClient
		// st.serverConnection, err = net.DialTCP("tcp", nil, tcpAddr)
		lib.PrintlnDebug("Dialed", st.rpcClient)
		if err != nil {
			lib.PrintlnError("Dial error", st.rpcClient, err)
			time.Sleep(200 * time.Millisecond)
			//shared.ErrorHandler(shared.GetFunction(), err.Error())
		} else {
			break
		}
	}
	lib.PrintlnDebug("Connected", st.rpcClient)
	// if addr != shared.NAMING_HOST+":"+shared.NAMING_PORT && shared.LocalAddr == "" {
	// 	//lib.PrintlnDebug("crhInfo.Conns[addr].LocalAddr().String()", crhInfo.Conns[addr].LocalAddr().String())
	// 	shared.LocalAddr = st.rpcClient.LocalAddr().String()
	// 	lib.PrintlnDebug("Got local addr", st.rpcClient)
	// }
}

func (st *RPC) CloseConnection() {
	err := st.rpcClient.Close()
	if err != nil {
		lib.PrintlnError(err)
	}
}

func (st *RPC) ReadString() (message string) {
	//TODO implement me
	panic("implement me")
}

func (st *RPC) WriteString(message string) {
	//TODO implement me
	panic("implement me")
}

func (st *RPC) Receive() ([]byte, error) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "RPC.Receive")
	msgFromServer := <-st.msgChan
	// sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	// // receive reply's size
	// _, err := st.serverConnection.Read(sizeOfMsgSize)
	// if err != nil {
	// 	lib.PrintlnError(shared.GetFunction(), err)
	// 	//shared.ErrorHandler(shared.GetFunction(), err.Error())
	// 	return nil, err
	// }
	// lib.PrintlnInfo("----------------------------------------->", shared.GetFunction(), "RPC read size")
	// // receive reply
	// msgFromServer := make([]byte, binary.LittleEndian.Uint32(sizeOfMsgSize), shared.NUM_MAX_MESSAGE_BYTES)
	// _, err = st.serverConnection.Read(msgFromServer)
	// if err != nil {
	// 	lib.PrintlnError(shared.GetFunction(), err)
	// 	//shared.ErrorHandler(shared.GetFunction(), err.Error())
	// 	return nil, err
	// }
	// lib.PrintlnInfo("----------------------------------------->", shared.GetFunction(), "RPC read message")
	return msgFromServer, nil
}

func (st *RPC) Call(serviceMethod string, args any, reply any) (err error) {
	c := make(chan error, 1)
	go func() { c <- st.rpcClient.Call(serviceMethod, args, reply) }()
	select {
	case err = <-c:
		return err
	case <-time.After(2 * time.Second):
		return errors.New("RPC Call timeout")
	}
}

func (st *RPC) Send(msgToServer []byte) error {
	// lib.PrintlnInfo("CRHRPC Version Not adapted")
	//sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE) // TODO dcruzb: create attribute to avoid doing this everytime

	// The message received from the server
	var msgFromServer []byte // := make([]byte, binary.LittleEndian.Uint32(sizeOfMsgSize), shared.NUM_MAX_MESSAGE_BYTES)
	err := st.Call("RPCRequest.Request", msgToServer, &msgFromServer)
	if err != nil {
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return err
	}

	// lib.PrintlnInfo("Got message from server")
	go func() {
		st.msgChan <- msgFromServer
	}()
	// lib.PrintlnInfo("Put message in msgChan")
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

type RPCRequest struct {
	msgChan   chan []byte
	replyChan chan []byte
}

func (rq RPCRequest) Request(request []byte, reply *[]byte) error {
	// lib.PrintlnInfo("Received message")
	go func() {
		rq.msgChan <- request
	}()
	// lib.PrintlnInfo("Forwarded message")
	replyMsg := <-rq.replyChan
	// lib.PrintlnInfo("Received reply")
	*reply = replyMsg
	return nil
}
