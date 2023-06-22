package messages

import (
	"net"

	"github.com/quic-go/quic-go"
)

type CRHInfo struct {
	EndPoint    EndPoint
	Conns       map[string]net.Conn
	QuicConns   map[string]quic.Connection
	QuicStreams map[string]quic.Stream
}
