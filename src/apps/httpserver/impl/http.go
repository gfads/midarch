package impl

import (
	"fmt"
	"gmidarch/development/messages"
)

func RequestListener(request messages.HttpRequest, response *messages.HttpResponse) {
	fmt.Println("Http.RequestListener request:", request.Method, request.Route)

	response.Protocol = "HTTP/1.1"
	response.Status = "200"
	response.Header.Fields = make(map[string]string)
	response.Header.Fields["content-type"] = "text/html; charset=UTF-8"
	response.Header.Fields["date"] = "Sun 06 Sep 2020 14:39:09 GMT"
	response.Body = "<html><h1>RequestListener Test Ok</h1></html>"
}
