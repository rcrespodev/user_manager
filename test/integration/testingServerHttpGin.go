package integration

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/api/v1/handlers/jwtAuth"
	"github.com/rcrespodev/user_manager/api/v1/routes"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/auth/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel"
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
	engine.Use(jwtAuth.ValidateJwt()) //Jwt Auth

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

	request, err := http.NewRequest(method, fmt.Sprintf("http://0.0.0.0:8080%v", path), bytes.NewReader(cmd.BodyRequest))
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Content-type", "application/json")

	if request.URL.Path != "/check-status" {
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
