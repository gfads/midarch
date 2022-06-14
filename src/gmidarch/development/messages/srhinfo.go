package messages

import (
	"log"
	"net"
)

type SRHInfo struct {
	EndPoint      	EndPoint               	// host, port
	Ln            	net.Listener           	// Listener
	Conns         	[]net.Conn             	// Set of connections
	CurrentConn		net.Conn               	// Current connection
	UDPConnection	*net.UDPConn			// UDP Connection
	RcvedMessages 	chan ReceivedMessages	// Buffer of messages received by the server
	Clients			[]*Client				// Connection Pool, possible connected
	Counter			int
	ExecuteForever	*bool
}

type ReceivedMessages struct {
	ToAddress 	string
	Chn 		net.Conn
	Msg 		[]byte
}

type Client struct {
	Ip				string
	Connection		net.Conn
	UDPConnection 	*net.UDPConn
	AdaptId			int
}

func (i SRHInfo) GetClientFromAddr(addr string, clients []*Client) *Client {
	for _, client := range clients {
		if client.Ip == addr {
			return client
		}
		log.Println("client without connection from the ip:", addr, " connection:", client.Connection)
	}

	return nil
}