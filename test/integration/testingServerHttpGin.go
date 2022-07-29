package integration

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/api/v1/endpoints"
	"github.com/rcrespodev/user_manager/api/v1/handlers/jwtAuth"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/auth/domain"
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
	engine.Use(jwtAuth.ValidateJwt()) //Jwt Auth

	for path, endpointData := range endPointsMap {
		engine.Handle(endpointData.HttpMethod, path, endpointData.Handler)
	}

	return &TestServerHttpGin{
		engine:    engine,
		endpoints: endPointsMap,
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

	endpointData, ok := t.endpoints[cmd.RelativePath]
	if ok {
		method = endpointData.HttpMethod
		path = cmd.RelativePath
	}

	request, err := http.NewRequest(method, fmt.Sprintf("http://0.0.0.0:8080%v", path), bytes.NewReader(cmd.BodyRequest))
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Content-type", "application/json")

	if request.URL.Path != endpoints.EndpointCheckStatus {
		jwt, _ := jwtDomain.SignJwt(uuid.New(), kernel.Instance.JwtConfig())
		request.Header.Set("Authorization", jwt)
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
