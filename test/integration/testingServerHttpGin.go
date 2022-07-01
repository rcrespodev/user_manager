package integration

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/rcrespodev/user_manager/api/v1/routes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
)

type TestServerHttpGin struct {
	engine *gin.Engine
	routes *routes.Routes
}

func NewTestServerHttpGin(routes *routes.Routes) *TestServerHttpGin {
	engine := gin.Default()
	for _, route := range routes.Routes {
		engine.Handle(route.HttpMethod, route.RelativePath, route.Handler)
	}
	return &TestServerHttpGin{
		engine: engine,
		routes: routes,
	}
}

type DoRequestCommand struct {
	BodyRequest  []byte
	RelativePath string
}

type Response struct {
	Header   http.Header
	HttpCode int
	Body     []byte
}

func (t TestServerHttpGin) DoRequest(cmd DoRequestCommand) Response {
	var method string
	var path string
	for _, route := range t.routes.Routes {
		if route.RelativePath == cmd.RelativePath {
			method = route.HttpMethod
			path = route.RelativePath
			break
		}
	}

	request, err := http.NewRequest(method, path, bytes.NewReader(cmd.BodyRequest))
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Content-type", "application/json")
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
