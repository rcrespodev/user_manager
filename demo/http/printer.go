package http

import (
	"encoding/json"
	"fmt"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"log"
	"net/http"
)

type Printer struct {
}

func (p *Printer) PrintStep(step int) {
	switch step {
	case 0:
		fmt.Printf("\n== 0. Setting data...\n")
	case 1:
		fmt.Printf("\n== 1. Setp one: Check Health App\n")
	case 2:
		fmt.Printf("\n== 2. Setp two: Register New User\n")
	case 3:
		fmt.Printf("\n== 3. Setp tree: Login User\n")
	case 4:
		fmt.Printf("\n== 4. Setp four: Get User Data\n")
	case 5:
		fmt.Printf("\n== 5. Setp Five: Login out user\n")
	case 6:
		fmt.Printf("\n== 6. Setp Six: Login user again\n")
	case 7:
		fmt.Printf("\n== 7. Setp Seven: Delete user\n")
	}
}

func (p *Printer) PrintRequestHeader(request *http.Request) {
	fmt.Printf("- Request:\n")
	fmt.Printf("\t- Making request to: %s%s%s\n", request.URL.Host, request.URL.Path, request.URL.RawQuery)
	fmt.Printf("\t- Method: %s\n", request.Method)
	if auth := request.Header.Get("Authorization"); auth != "" {
		fmt.Printf("\t- Request Token: %s...\n", auth[1:80])
	}
}

func (p *Printer) PrintResponseHeader(response *http.Response) {
	fmt.Printf("- Response:\n")
	fmt.Printf("\t- Status code: %d\n", response.StatusCode)
	token := response.Header.Get("Token")
	if token != "" {
		fmt.Printf("\t- JWT: %s...\n", token[1:80])
	}
}

func (p *Printer) PrintResponseCommand(response *api.CommandResponse) {
	fmt.Printf("\t- Response body:\n")
	p.PrintResponseMessage(response.Message)
}

func (p *Printer) PrintResponseQuery(response *api.QueryResponse) {
	fmt.Printf("\t- response body:\n")
	p.PrintResponseMessage(response.Message)
	p.PrintResponseData(response.Data)
}

func (p *Printer) PrintResponseMessage(message message.MessageData) {
	fmt.Printf("\t\t- Message:\n")
	fmt.Printf("\t\t\t- ObjectId: %v\n", message.ObjectId)
	fmt.Printf("\t\t\t- MessageId: %v\n", message.MessageId)
	fmt.Printf("\t\t\t- MessagePkg: %v\n", message.MessagePkg)
	fmt.Printf("\t\t\t- Variables: %v\n", message.Variables)
	fmt.Printf("\t\t\t- Text: %v\n", message.Text)
	fmt.Printf("\t\t\t- Time: %v\n", message.Time)
	fmt.Printf("\t\t\t- ClientError: %v\n", message.ClientErrorType)
}

func (p *Printer) PrintResponseData(data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\t\t- Data:\n")
	fmt.Printf("\t\t\t%s\n", string(b))
}

func (p *Printer) PrintRequestBody(body []byte) {
	fmt.Printf("\t- Request Body:\n")
	fmt.Printf("\t\t%s\n", string(body))
}
