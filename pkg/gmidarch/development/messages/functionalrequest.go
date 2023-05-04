package messages

type FunctionalRequest struct {
	Op     string        `json:"Op"`
	Params []interface{} `json:"Params"`
}
