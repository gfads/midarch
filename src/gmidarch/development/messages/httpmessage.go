package messages

type HttpMessage struct {
	Type string // Request | Response
	Method string
	Route string
	QueryParameters string
	Protocol string
	Status int
	Headers Header
	Body string
}

type Header struct {
	Fields map[string]string
}