package http

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/api/v1/endpoints"
	delete "github.com/rcrespodev/user_manager/pkg/app/user/application/commands/delete"
	"github.com/rcrespodev/user_manager/pkg/app/user/application/commands/login"
	"github.com/rcrespodev/user_manager/pkg/app/user/application/commands/register"
	"log"
	"net/http"
	"time"
)

const (
	userAlias      = "user_demo_alias"
	userName       = "user_demo_name"
	userSecondName = "user_demo_second_name"
	userEmail      = "user_demo_email@gmail.com"
	userPassword   = "Linux64bits$"
)

type Demo struct {
	clientHttp *ClientHttp
	printer    *Printer
	token      string
	userUuid   string
}

func NewDemo(clientHttp *ClientHttp) *Demo {
	return &Demo{
		clientHttp: clientHttp,
	}
}

func (d *Demo) Exec() {
	d.checkHealth(1)
	time.Sleep(time.Second * 1)
	d.registerNewUser(2)
	time.Sleep(time.Second * 1)
	d.loginUser(3)
	time.Sleep(time.Second * 1)
	d.getUserInfo(4)
	time.Sleep(time.Second * 1)
	d.logOutUser(5)
	time.Sleep(time.Second * 1)
	d.loginUser(6)
	time.Sleep(time.Second * 1)
	d.deleteUser(7)
}

func (d *Demo) checkHealth(step int) {
	d.printer.PrintStep(step)
	header, body := d.clientHttp.DoRequest(DoRequestCommand{
		RelPath: endpoints.EndpointCheckStatus,
		Method:  http.MethodGet,
	})
	d.printer.PrintResponseHeader(header)
	d.printer.PrintResponseQuery(d.parseResponseQuery(body))
}

func (d *Demo) registerNewUser(step int) {
	d.printer.PrintStep(step)

	d.userUuid = uuid.New().String()
	cmd := register.ClientArgs{
		Uuid:       d.userUuid,
		Alias:      userAlias,
		Name:       userName,
		SecondName: userSecondName,
		Email:      userEmail,
		Password:   userPassword,
	}

	bodyRequest, err := json.Marshal(cmd)
	if err != nil {
		log.Fatal(err)
	}

	response, body := d.clientHttp.DoRequest(DoRequestCommand{
		RelPath: endpoints.EndpointUser,
		Method:  http.MethodPost,
		Body:    bodyRequest,
	})
	d.setToken(response)
	d.printer.PrintResponseHeader(response)
	d.printer.PrintResponseCommand(d.parseResponseCommand(body))
}

func (d *Demo) parseResponseCommand(respBody []byte) *api.CommandResponse {
	var gotRespBody *api.CommandResponse
	if err := json.Unmarshal(respBody, &gotRespBody); err != nil {
		log.Panicln(err)
	}
	return gotRespBody
}

func (d *Demo) parseResponseQuery(respBody []byte) *api.QueryResponse {
	var gotRespBody *api.QueryResponse
	if err := json.Unmarshal(respBody, &gotRespBody); err != nil {
		log.Panicln(err)
	}
	return gotRespBody
}

func (d *Demo) getUserInfo(step int) {
	d.printer.PrintStep(step)

	response, body := d.clientHttp.DoRequest(DoRequestCommand{
		RelPath:       fmt.Sprintf("%s?email=%s", endpoints.EndpointUser, userEmail),
		Method:        http.MethodGet,
		Authorization: d.token,
		Body:          nil,
	})
	d.setToken(response)
	d.printer.PrintResponseHeader(response)
	d.printer.PrintResponseQuery(d.parseResponseQuery(body))
}

func (d *Demo) logOutUser(step int) {
	d.printer.PrintStep(step)

	response, body := d.clientHttp.DoRequest(DoRequestCommand{
		RelPath:       endpoints.EndpointUserLogin,
		Method:        http.MethodDelete,
		Authorization: d.token,
		Body:          nil,
	})
	d.setToken(response)
	d.printer.PrintResponseHeader(response)
	d.printer.PrintResponseCommand(d.parseResponseCommand(body))
}

func (d *Demo) loginUser(step int) {
	d.printer.PrintStep(step)

	cmd := login.ClientArgs{
		AliasOrEmail: userEmail,
		Password:     userPassword,
	}

	bodyRequest, err := json.Marshal(cmd)
	if err != nil {
		log.Fatal(err)
	}

	response, body := d.clientHttp.DoRequest(DoRequestCommand{
		RelPath: endpoints.EndpointUserLogin,
		Method:  http.MethodPost,
		Body:    bodyRequest,
	})
	d.setToken(response)
	d.printer.PrintResponseHeader(response)
	d.printer.PrintResponseCommand(d.parseResponseCommand(body))
}

func (d *Demo) deleteUser(step int) {
	d.printer.PrintStep(step)

	cmd := delete.ClientArgs{UserUuid: d.userUuid}

	bodyRequest, err := json.Marshal(cmd)
	if err != nil {
		log.Fatal(err)
	}

	response, body := d.clientHttp.DoRequest(DoRequestCommand{
		RelPath:       endpoints.EndpointUser,
		Method:        http.MethodDelete,
		Authorization: d.token,
		Body:          bodyRequest,
	})
	d.setToken(response)
	d.printer.PrintResponseHeader(response)
	d.printer.PrintResponseCommand(d.parseResponseCommand(body))
}

func (d *Demo) setToken(response *http.Response) {
	if token := response.Header.Get("Token"); token != "" {
		d.token = token
	}
}
