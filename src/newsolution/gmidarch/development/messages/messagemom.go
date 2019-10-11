package messages

type MessageMOM struct{
	Header Header
	PayLoad string
}

type Header struct {
	Destination string
}

