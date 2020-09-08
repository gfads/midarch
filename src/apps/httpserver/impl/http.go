package impl

import (
	"apps/fibomiddleware/impl"
	"fmt"
	"gmidarch/development/messages"
	"strconv"
	"strings"
)

func RequestListener(request messages.HttpRequest, response *messages.HttpResponse) {
	fmt.Println("Http.RequestListener request:", request.Method, request.Route, request.QueryParameters)

	response.Protocol = "HTTP/1.1"
	response.Header.Fields = make(map[string]string)
	response.Header.Fields["date"] = "Sun 06 Sep 2020 14:39:09 GMT"

	switch request.Route {
	case "/test":
		response.Status = "200"
		response.Header.Fields["content-type"] = "text/html; charset=UTF-8"

		if request.QueryParameters != "" {
			fmt.Println("Http.RequestListener queryParameters:", request.QueryParameters)
			parameters := strings.Split(request.QueryParameters, "&")
			response.Body = "<html><h1>RequestListener Test Ok</h1>"
			response.Body += "<ul>"
			for _, parameter := range parameters {
				fmt.Println("Http.RequestListener parameter:", parameter)
				response.Body += "<li>"+parameter
			}
			response.Body += "</ul></html>"
		} else {
			response.Body = "<html><h1>RequestListener Test Ok (without query parameters)</h1></html>"
		}
	case "/Fibo":
		response.Status = "200"
		response.Header.Fields["content-type"] = "text/html; charset=UTF-8"

		if request.QueryParameters != "" {
			fmt.Println("Http.RequestListener queryParameters:", request.QueryParameters)
			_p, _ := strconv.Atoi(strings.Split(request.QueryParameters, "=")[1])

			_r := impl.Fibonacci{}.F(_p)

			response.Body = strconv.Itoa(_r)
		} else {
			response.Status = "400"
		}
	default:
		response.Status = "400"
	}
}
