package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ClientHttp struct {
	host     string
	port     string
	basePath string
	client   *http.Client
	printer  *Printer
}

type DoRequestCommand struct {
	RelPath       string
	Method        string
	Body          []byte
	Authorization string
}

func NewClientHttp(host, port string) *ClientHttp {
	return &ClientHttp{
		host:     host,
		port:     port,
		basePath: fmt.Sprintf("http://%s:%s", host, port),
		client: &http.Client{
			Timeout: 200 * time.Second,
		},
	}
}

func (c *ClientHttp) DoRequest(command DoRequestCommand) (header *http.Response, body []byte) {
	fullPath := fmt.Sprintf("%s%s", c.basePath, command.RelPath)
	reader := bytes.NewReader(command.Body)
	request, err := http.NewRequest(command.Method, fullPath, reader)
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Content-type", "application/json")
	if command.Authorization != "" {
		request.Header.Set("Authorization", command.Authorization)
	}
	response, err := c.client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	c.printer.PrintRequestHeader(request)
	if command.Body != nil {
		c.printer.PrintRequestBody(command.Body)
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return response, respBody
}
