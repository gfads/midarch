package messages

type MessageMOM struct{
	Header MOMHeader
	Payload interface{}
}

type MOMHeader struct {
	Destination string
}
