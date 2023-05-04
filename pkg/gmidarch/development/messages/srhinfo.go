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

func (c Client) Initialize() {
	c.Ip = ""
	if c.Connection != nil {
		c.Connection.Close()
	}
	c.Connection = nil
	if c.UDPConnection != nil {
		c.UDPConnection.Close()
	}
	c.UDPConnection = nil
	c.AdaptId = 0
}

func (i SRHInfo) GetClientFromAddr(addr string, clients []*Client) *Client {
	for _, client := range clients {
		if client.Ip == addr {
			return client
		}
	}
	log.Println("IP without client from the ip:", addr)

	return nil
}

func (i SRHInfo) GetClientFromAdaptId(adaptId int, clients []*Client) *Client {
	for _, client := range clients {
		if client.AdaptId == adaptId {
			return client
		}
	}
	log.Println("AdaptId without client => adaptId:", adaptId)

	return nil
}

func Remove(slice []*Client, idx int) []*Client {
	var newSlice []*Client

	if len(slice) == idx+1 {
		newSlice = append(slice[:idx])
	} else {
		newSlice = append(slice[:idx], slice[idx+1:]...)
	}

	return newSlice
}