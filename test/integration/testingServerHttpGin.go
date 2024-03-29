package integration

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rcrespodev/user_manager/api/v1/endpoints"
	"github.com/rcrespodev/user_manager/api/v1/handlers/ginMiddleware"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
)

type TestServerHttpGin struct {
	engine    *gin.Engine
	endpoints endpoints.Endpoints
}

func NewTestServerHttpGin(endPointsMap endpoints.Endpoints) *TestServerHttpGin {
	engine := gin.Default()
	engine.Use(ginMiddleware.MiddlewareHandlerFunc()) //Prd middleware

	for _, endpointData := range endPointsMap {
		engine.Handle(endpointData.HttpMethod, endpointData.RelPath, endpointData.Handler)
	}

	return &TestServerHttpGin{
		engine:    engine,
		endpoints: endPointsMap,
	}
}

type DoRequestCommand struct {
	BodyRequest  []byte
	RelativePath string
	Method       string
	Uuid         string
	Token        string
	QueryString  string
}

type Response struct {
	Header   http.Header
	HttpCode int
	Body     []byte
}

func (t TestServerHttpGin) DoRequest(cmd DoRequestCommand) Response {
	var method string
	var path string

	endpointData, ok := t.endpoints[endpoints.BuildEndpointKey(cmd.RelativePath, cmd.Method)]
	if ok {
		method = endpointData.HttpMethod
		path = endpointData.RelPath
	}

	if cmd.QueryString != "" {
		path = fmt.Sprintf("%s%s", cmd.RelativePath, cmd.QueryString)
	}

	request, err := http.NewRequest(method, fmt.Sprintf("http://0.0.0.0:8080%v", path), bytes.NewReader(cmd.BodyRequest))
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Content-type", "application/json")

	if endpointData.AuthValidation {
		token := cmd.Token
		if token == "" {
			token, err = kernel.Instance.Jwt().CreateNewToken(cmd.Uuid)
			if err != nil {
				log.Fatal(err)
			}
		}
		request.Header.Set("Authorization", token)
	}

	writer := httptest.NewRecorder()
	t.engine.ServeHTTP(writer, request)

	bodyResponse, err := ioutil.ReadAll(writer.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return Response{
		Header:   writer.Header(),
		HttpCode: writer.Code,
		Body:     bodyResponse,
	}
}
