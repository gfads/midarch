package miop

type Packet struct {
	Hdr Header
	Bd Body
}

type Header struct{
	Magic string
	Version string
	ByteOrder bool
	MessageType int
	Size int
}

type Body struct {
	ReqHeader RequestHeader
	ReqBody RequestBody
	RepHeader ReplyHeader
	RepBody ReplyBody
}

type RequestHeader struct {
	Context string
	RequestId int
	ResponseExpected bool
	Key int
	Operation string
}

type RequestBody struct {
	Body []interface{}
}

type ReplyHeader struct {
	Context string
	RequestId int
	Status int
}

type ReplyBody struct {
	OperationResult [] interface{}
}


