package messages

type RequestHeader struct{
	Context string
	RequestId int
	ResponseExpected bool
	Key int
	Operation string
}

type RequestBody struct {
	Args interface{}
}

type ReplyHeader struct {
	Context string
	RequestId int
	Status int
}

type ReplyBody struct {
	Args interface{}
}

type MIOP struct {
	Header MIOPHeader
	Body MIOPBody
}

type MIOPHeader struct {
	Magic string
	Version string
	ByteOrder bool
	MessageType int
	MessageSize int
}

type MIOPBody struct {
	RequestHeader RequestHeader
	RequestBody RequestBody
	ReplyHeader ReplyHeader
	ReplyBody ReplyBody
}