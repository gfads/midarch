package messages

import "net/http"

type HttpMessage struct {
	Response http.ResponseWriter
	Request *http.Request
}