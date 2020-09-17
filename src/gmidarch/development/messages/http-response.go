package messages

import (
	"strings"
)

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
	resp.Header.Fields = make(map[string]string)
	lines := strings.Split(payload, "\n")
	//fmt.Println("HttpResponse.Unmarshal payload:", payload)

	bodyStarted := false
	for _, line := range lines {
		if resp.Protocol == "" {
			startLine := strings.Fields(line)
			//fmt.Println("HttpResponse.Unmarshal startLine:", startLine)
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