package messages

type SAMRequest struct {
	Op     string        `json:"Op"`
	Params []interface{} `json:"Params"`
}
