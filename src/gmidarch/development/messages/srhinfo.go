package messages

import (
	"net"
)

type SRHInfo struct {
	EndPoint      EndPoint               // host, port
	Ln            net.Listener           // Listener
	Conns         []net.Conn             // Set of connections
	CurrentConn   net.Conn               // Current connection
	RcvedMessages chan ReceivedMessages  // Buffer of messages received by the server
}

type ReceivedMessages struct {
	Chn net.Conn
	Msg []byte
}
