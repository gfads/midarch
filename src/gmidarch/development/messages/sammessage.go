package messages

// Software architecture message
type SAMessage struct {
	From string   `json:"From"`
	To string     `json:"To"`
	Payload interface{} `json:"Payload"`
}