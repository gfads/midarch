package miop

type MiopPacket struct {
	Hdr Header `json:"header"`
	Bd  Body   `json:"body"`
}

type Header struct {
	Magic       string `json:"magic"`
	Version     string `json:"version"`
	ByteOrder   bool   `json:"byteOrder"`
	MessageType int    `json:"messageType"`
	Size        int    `json:"size"`
}

type Body struct {
	ReqHeader RequestHeader `json:"reqHead"`
	ReqBody   RequestBody   `json:"reqBody"`
	RepHeader ReplyHeader   `json:"repHeader"`
	RepBody   ReplyBody     `json:"repBody"`
}

type RequestHeader struct {
	Context          string `json:"context"`
	RequestId        int    `json:"requestId"`
	ResponseExpected bool   `json:"responseExpected"`
	Key              int    `json:"key"`
	Operation        string `json:"operation"`
}

type RequestBody struct {
	Body []interface{} `json:"body"`
}

type ReplyHeader struct {
	Context   string `json:"context"`
	RequestId int    `json:"requestId"`
	Status    int    `json:"status"`
}

type ReplyBody struct {
	OperationResult interface{} `json:"operationresult"`
}

func CreateReqPacket(op string, params []interface{}, adaptId int) MiopPacket {

	// MIOP Header
	hd := Header{}
	hd.Magic = "MIOP"
	hd.ByteOrder = true
	hd.MessageType = 1
	hd.Size = 1024
	hd.Version = "version 1.0"

	// MIOP Body -- Request Header
	bd := Body{}
	reqHd := RequestHeader{}
	reqHd.Operation = op
	reqHd.RequestId = adaptId
	reqHd.Context = "context"
	reqHd.Key = 1313
	reqHd.ResponseExpected = true
	bd.ReqHeader = reqHd

	// Request Body
	reqBd := RequestBody{}
	reqBd.Body = params
	bd.ReqBody = reqBd

	// Body
	r := MiopPacket{Hdr: hd, Bd: bd}

	return r
}

func CreateRepPacket(params interface{}) MiopPacket {

	// MIOP Header
	hd := Header{}
	hd.Magic = "MIOP"
	hd.ByteOrder = true
	hd.MessageType = 1
	hd.Size = 1024
	hd.Version = "version 1.0"

	// MIOP Body -- Reply Header
	bd := Body{}
	repHd := ReplyHeader{}
	repHd.Context = "context"
	repHd.Status = 0
	repHd.RequestId = 0
	bd.RepHeader = repHd

	// Reply Body
	repBd := ReplyBody{}
	repBd.OperationResult = params
	bd.RepBody = repBd

	// Body
	r := MiopPacket{Hdr: hd, Bd: bd}

	return r
}
