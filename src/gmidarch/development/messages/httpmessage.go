package messages

import (
	"fmt"
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
	//fmt.Println("HttpResponse.Unmarshal lines:", lines)
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

func (resp *HttpResponse) Unmarshal(payload string) {
	//fmt.Println("HttpResponse.Unmarshal payload:", payload)
	resp.Header.Fields = make(map[string]string)
	lines := strings.Split(payload, "\n")
	bodyStarted := false
	for _, line := range lines {
		if resp.Protocol == "" {
			startLine := strings.Fields(line)
			fmt.Println("HttpResponse.Unmarshal startLine:", startLine)
			resp.Protocol = startLine[0]
			resp.Status = startLine[1]

			continue
		}

		if strings.TrimSpace(line) == "" {
			bodyStarted = true
			continue
		}

		if !bodyStarted {
			header := strings.Split(line, ": ")
			//fmt.Println("HttpResponse.Unmarshal header:", header)
			resp.Header.Fields[header[0]] = header[1]
		}else{
			resp.Body += line
		}
	}
}