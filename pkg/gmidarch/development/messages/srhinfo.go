package messages

import (
	"log"
	"net"

	"github.com/quic-go/quic-go"
)

type SRHInfo struct {
	EndPoint       EndPoint              // host, port
	Ln             net.Listener          // Listener
	QUICLn         quic.Listener         // Listener
	Conns          []net.Conn            // Set of connections
	QUICConns      []quic.Connection     // Set of connections
	CurrentConn    net.Conn              // Current connection
	UDPConnection  *net.UDPConn          // UDP Connection
	RcvedMessages  chan ReceivedMessages // Buffer of messages received by the server
	Clients        []*Client             // Connection Pool, possible connected
	Counter        int
	ExecuteForever *bool
}

type ReceivedMessages struct {
	ToAddress  string
	Conn       net.Conn
	QUICStream quic.Stream
	Msg        []byte
}

type Client struct {
	Ip             string
	Connection     net.Conn
	UDPConnection  *net.UDPConn
	QUICConnection quic.Connection
	QUICStream     quic.Stream
	AdaptId        int
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
	if c.QUICStream != nil {
		c.QUICStream.Close()
	}
	c.QUICStream = nil
	if c.QUICConnection != nil {
		c.QUICConnection.CloseWithError(0, "Initialized SRH")
	}
	c.QUICConnection = nil

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
