package messages

// Software architecture message
type SAMessage struct {
	From string   		`json:"From"`
	To string     		`json:"To"`
	ToAddr string		`json:"ToAddr"`
	Payload interface{} `json:"Payload"`
}