package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/api/v1/endpoints"
	"github.com/rcrespodev/user_manager/pkg/app/user/application/commands/register"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Printf("1. setp one\n")
	checkHealth()
	fmt.Printf("2. setp two\n")
	registerUser()
}

func registerUser() {
	relPath := endpoints.EndpointRegisterUser
	fmt.Printf("register new user...\n")
	fmt.Printf("\trequest to: %s\n", relPath)

	random := 27.2
	//random := rand.Float32()
	cmd := register.ClientArgs{
		Uuid:       uuid.New().String(),
		Alias:      fmt.Sprintf("foo_alias_%.0f", random),
		Name:       fmt.Sprintf("foo_name%.0f", random),
		SecondName: fmt.Sprintf("foo_second%.0f", random),
		Email:      fmt.Sprintf("foo_email_%.0f@gmail.com", random),
		Password:   "Linux64bits$",
	}
	bytesCmd, err := json.Marshal(cmd)
	if err != nil {
		log.Panic(err)
	}
	doRequest(http.MethodPost, relPath, bytesCmd)
}

func checkHealth() {
	relPath := endpoints.EndpointCheckStatus
	fmt.Printf("cheking App Status...\n")
	fmt.Printf("\trequest to: %s\n", relPath)
	doRequest(http.MethodGet, relPath, nil)
}

func doRequest(method, relPath string, body []byte) {
	fullPath := fmt.Sprintf("http://0.0.0.0:8080%v", relPath)
	request, err := http.NewRequest(method, fullPath, bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Content-type", "application/json")

	client := &http.Client{
		Timeout: 200 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\t- status code: %d\n", response.StatusCode)

	switch method {
	case http.MethodPost:
		var gotRespBody *api.CommandResponse
		if err := json.Unmarshal(respBody, &gotRespBody); err != nil {
			log.Panicln(err)
		}
		fmt.Printf("\t- body:\n\t\tmessage: %+v\n", gotRespBody.Message)
	case http.MethodGet:
		var gotRespBody *api.QueryResponse
		if err := json.Unmarshal(respBody, &gotRespBody); err != nil {
			log.Panicln(err)
		}
		fmt.Printf("\t- body:\n\t\tmessage: %+v\n\t\tdata: %+v\n", gotRespBody.Message, gotRespBody.Data)
	}
}
