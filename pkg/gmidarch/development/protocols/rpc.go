package protocols

import (
	"log"
	"net"
	"net/rpc"

	"github.com/gfads/midarch/pkg/gmidarch/development/generic"
)

type Request struct {
	Message []byte
}

type Reply struct {
	Message []byte
}

type RPC struct {
	ip        string
	port      string
	listener  net.Listener
	rpcClient *rpc.Client
	msgChan   chan []byte
}

func (this *RPC) StartServer(ip, port string, initialConnections int) {
	//TODO implement me
	panic("implement me")
}

func (this *RPC) StopServer() {
	//TODO implement me
	panic("implement me")
}

func (this *RPC) AvailableConnectionFromPool() (available bool, idx int) {
	//TODO implement me
	return true, 0
}

func (this *RPC) WaitForConnection(cliIdx int) (cl *generic.Client) { // TODO if cliIdx >= inicitalConnections => need to append to the slice
	//TODO implement me
	panic("implement me")
}

func (this *RPC) ConnectToServer(ip, port string) {
	this.msgChan = make(chan []byte)

	// connect to server
	rpcClient, err := rpc.DialHTTP("tcp", ip+":"+port)
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	this.rpcClient = rpcClient
}

func (this *RPC) CloseConnection() {
	//TODO implement me
	panic("implement me")
}

func (this *RPC) ReadString() string {
	//TODO implement me
	panic("implement me")
}

func (this *RPC) WriteString(message string) {
	//TODO implement me
	panic("implement me")
}

func (this *RPC) Receive() ([]byte, error) {
	msg := <-this.msgChan
	return msg, nil
}

func (this *RPC) Send(msgToServer []byte) error {
	reply := Reply{}
	err := this.rpcClient.Call("call", Request{Message: msgToServer}, &reply)
	if err != nil {
		return err
	}
	this.msgChan <- reply.Message
	return nil
}

func (this *RPC) GetClients() (clients []*generic.Client) {
	//TODO implement me
	panic("implement me")
}

func (this *RPC) GetClient(idx int) (client generic.Client) {
	//TODO implement me
	panic("implement me")
}

func (this *RPC) GetClientFromAddr(addr string) (client generic.Client) {
	//TODO implement me
	panic("implement me")
}

func (this *RPC) AddClient(client generic.Client, idx int) {
	//TODO implement me
	panic("implement me")
}

func (this *RPC) ResetClients() {
	//TODO implement me
	panic("implement me")
}
