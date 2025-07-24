package docker

type RemoteOperationFactor int

const (
	Fibonacci RemoteOperationFactor = iota
	SendFile
)

func (kind RemoteOperationFactor) toString() string {
	switch kind {
	case Fibonacci:
		return "Fibonacci"
	case SendFile:
		return "SendFile"
	}
	panic("Kind conversion to string using unlisted kind")
}

type TransportProtocolFactor int

const (
	Udp TransportProtocolFactor = iota
	Tcp
	Tls
	Rpc
	Quic
	Http
	Https
	Http2
	E_Rpc
	E_Grpc
	E_Rmq
	UdpTcp
	TcpTls
	RpcQuic
	RpcHttp
	QuicHttp2
	TcpHttp
	TlsHttp2
)

func (kind TransportProtocolFactor) IsEvolutive() bool {
	switch kind {
	case UdpTcp, TcpTls, RpcQuic, RpcHttp, QuicHttp2, TcpHttp, TlsHttp2:
		return true
	}
	return false
}

func (kind TransportProtocolFactor) IsMidArch() bool {
	return !(kind == E_Rpc || kind == E_Grpc || kind == E_Rmq)
}

func (kind TransportProtocolFactor) getEvolutiveProtocols() (TransportProtocolFactor, TransportProtocolFactor) {
	switch kind {
	case UdpTcp:
		return Udp, Tcp
	case TcpTls:
		return Tcp, Tls
	case RpcQuic:
		return Rpc, Quic
	case RpcHttp:
		return Rpc, Http
	case QuicHttp2:
		return Quic, Http2
	case TcpHttp:
		return Tcp, Http
	case TlsHttp2:
		return Tls, Http2
	}
	panic("Kind get protocols using unlisted kind")
}

func (kind TransportProtocolFactor) toString() string {
	switch kind {
	case Udp:
		return "UDP"
	case Tcp:
		return "TCP"
	case Tls:
		return "TLS"
	case Quic:
		return "QUIC"
	case Rpc:
		return "RPC"
	case Http:
		return "HTTP"
	case Https:
		return "HTTPS"
	case Http2:
		return "HTTP2"
	case E_Rpc:
		return "E_RPC"
	case E_Grpc:
		return "E_GRPC"
	case E_Rmq:
		return "E_RMQ"
	case UdpTcp:
		return "UdpTcp"
	case TcpTls:
		return "TcpTls"
	case RpcQuic:
		return "RpcQuic"
	case RpcHttp:
		return "RpcHttp"
	case QuicHttp2:
		return "QuicHttp2"
	case TcpHttp:
		return "TcpHttp"
	case TlsHttp2:
		return "TlsHttp2"
	}
	panic("Kind conversion to string using unlisted kind")
}

func (kind TransportProtocolFactor) getComponentName() string {
	switch kind {
	case Udp:
		return "srhudp"
	case Tcp:
		return "srhtcp"
	case Tls:
		return "srhtls"
	case Quic:
		return "srhquic"
	case Rpc:
		return "srhrpc"
	case Http:
		return "srhhttp"
	case Https:
		return "srhttps"
	case Http2:
		return "srhhttp2"
	}
	panic("Kind get component name using unlisted kind")
}

func (kind TransportProtocolFactor) createStackCommand() string {
	switch kind {
	case Udp:
		return "docker stack deploy -c ./evaluation/experiments_v14/docker/dc-fibonaccidistributed-udp.yml fibonaccidistributed-udp"
	case Tcp:
		return "docker stack deploy -c ./evaluation/experiments_v14/docker/dc-fibonaccidistributed-tcp.yml fibonaccidistributed-tcp"
	case Tls:
		return "docker stack deploy -c ./evaluation/experiments_v14/docker/dc-fibonaccidistributed-tls.yml fibonaccidistributed-tls"
	case Quic:
		return "docker stack deploy -c ./evaluation/experiments_v14/docker/dc-fibomiddleware-quic.yml fibomiddleware-quic"
	case Rpc:
		return "docker stack deploy -c ./evaluation/experiments_v14/docker/dc-fibomiddleware-rpc.yml fibomiddleware-rpc"
	case Http:
		return "docker stack deploy -c ./evaluation/experiments_v14/docker/dc-fibomiddleware-http.yml fibomiddleware-http"
	case Https:
		return "docker stack deploy -c ./evaluation/experiments_v14/docker/dc-fibomiddleware-https.yml fibomiddleware-https"
	case Http2:
		return "docker stack deploy -c ./evaluation/experiments_v14/docker/dc-fibomiddleware-http2.yml fibomiddleware-http2"
	case E_Rpc:
		return "docker stack deploy -c ./evaluation/experiments_v14/docker/dc-fiborpc.yml fiborpc"
	case E_Grpc:
		return "docker stack deploy -c ./evaluation/experiments_v14/docker/dc-fibogrpc.yml fibogrpc"
	case E_Rmq:
		return "docker stack deploy -c ./evaluation/experiments_v14/docker/dc-fibormq.yml fibormq"
	case UdpTcp:
		return "docker stack deploy -c ./evaluation/experiments_v14/docker/dc-fibonaccidistributed-udptcp.yml fibonaccidistributed-udptcp"
	}
	panic("Kind create stack command using unlisted kind")
}

func (kind TransportProtocolFactor) removeStackCommand() string {
	switch kind {
	case Udp:
		return "docker stack rm fibonaccidistributed-udp"
	case Tcp:
		return "docker stack rm fibonaccidistributed-tcp"
	case Tls:
		return "docker stack rm fibonaccidistributed-tls"
	case Quic:
		return "docker stack rm fibomiddleware-quic"
	case Rpc:
		return "docker stack rm fibomiddleware-rpc"
	case Http:
		return "docker stack rm fibomiddleware-http"
	case Https:
		return "docker stack rm fibomiddleware-https"
	case Http2:
		return "docker stack rm fibomiddleware-http2"
	case E_Rpc:
		return "docker stack rm fiborpc"
	case E_Grpc:
		return "docker stack rm fibogrpc"
	case E_Rmq:
		return "docker stack rm fibormq"
	case UdpTcp:
		return "docker stack rm fibonaccidistributed-udptcp"
	}
	panic("Kind remove stack command using unlisted kind")
}
