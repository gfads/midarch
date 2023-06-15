package messages

type ProxyInfo struct {
	In_Channel  chan SAMessage
	Out_Channel chan SAMessage
}
