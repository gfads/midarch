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

func (req *HttpRequest) Marshal() []byte {
	var request string
	if req.QueryParameters != "" {
		request = req.Method + " " + req.Route + "?" + req.QueryParameters + " " + req.Protocol + "\n"
	} else {
		request = req.Method + " " + req.Route +                             " " + req.Protocol + "\n"
	}

	for key, value := range req.Header.Fields {
		request += key + ": " + value + "\n"
	}
	request += "\n" + req.Body

	return []byte(request)
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
			continue
		}

		if !bodyStarted {
			header := strings.Split(line, ": ")
			req.Header.Fields[header[0]] = header[1]
		}else{
			req.Body += line
		}
	}
}