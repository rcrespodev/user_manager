package integration

import (
	"fmt"
	"io"
	"net/http"
)

type NewHttpRequestCommand struct {
	Method      string
	Host        string
	Port        string
	Path        string
	Body        io.Reader
	ContentType string
}

func NewHttpRequest(command NewHttpRequestCommand) *Response {
	var err error
	var response *http.Response

	endpoint := fmt.Sprintf("http://%v:%v%v", command.Host, command.Port, command.Path)

	if command.Method == http.MethodGet {
		if response, err = http.Get(endpoint); err != nil {
			return badRequest()
		}
	}
	if command.Method == http.MethodPost {
		if response, err = http.Post(endpoint, command.ContentType, command.Body); err != nil {
			return badRequest()
		}
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return badRequest()
	}

	defer response.Body.Close()

	return &Response{
		Header:   response.Header,
		HttpCode: response.StatusCode,
		Body:     body,
	}
}

func badRequest() *Response {
	return &Response{
		Header:   nil,
		HttpCode: 500,
		Body:     nil,
	}
}
