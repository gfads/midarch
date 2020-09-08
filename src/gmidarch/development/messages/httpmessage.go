package messages

import (
	"strings"
)

type Header struct {
	Fields map[string]string
}

type HttpRequest struct {
	Method string
	Route string
	QueryParameters string
	Protocol string
	Header Header
	Body string
}

func (req *HttpRequest) Unmarshal(payload []byte) {
	req.Header.Fields = make(map[string]string)
	lines := strings.Split(string(payload), "\n")
	bodyStarted := false
	for _, line := range lines {
		if req.Method == "" {
			startLine := strings.Fields(line)
			req.Method = startLine[0]
			req.Protocol = startLine[2]

			route := strings.Split(startLine[1], "?")
			req.Route = route[0]
			if len(route) > 1 {
				req.QueryParameters = route[1]
			}

			continue
		}

		if strings.TrimSpace(line) == "" {
			bodyStarted = true
		}

		if !bodyStarted {
			header := strings.Split(line, ": ")
			req.Header.Fields[header[0]] = header[1]
		}else{
			req.Body += line
		}
	}
}

type HttpResponse struct {
	Protocol string
	Status string
	Header Header
	Body string
}

func (resp HttpResponse) Marshal() []byte {
	response := resp.Protocol + " " + resp.Status + "\n"
	for key, value := range resp.Header.Fields {
		response += key + ": " + value + "\n"
	}
	response += "\n" + resp.Body

	return []byte(response)
}