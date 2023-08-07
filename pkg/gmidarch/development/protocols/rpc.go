package protocols

import (
	"log"
	"net"
	"net/rpc"
)

type RPC struct {
	ip        string
	port      string
	listener  net.Listener
	rpcClient *rpc.Client
}

func (RPC) StartServer(ip, port string) {
	//TODO implement me
	panic("implement me")
}

func (RPC) StopServer() {
	//TODO implement me
	panic("implement me")
}

func (RPC) WaitForConnection() {
	//TODO implement me
	panic("implement me")
}

func (this *RPC) ConnectToServer(ip, port string) {
	// connect to server
	rpcClient, err := rpc.DialHTTP("tcp", ip+":"+port)
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	this.rpcClient = rpcClient
}

func (RPC) CloseConnection() {
	//TODO implement me
	panic("implement me")
}

func (RPC) Send() {
	//TODO implement me
	panic("implement me")
}

func (RPC) Receive() {
	//TODO implement me
	panic("implement me")
}
