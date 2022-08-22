package main

import "github.com/rcrespodev/user_manager/demo/http"

const (
	host = "0.0.0.0"
	port = "8080"
)

func main() {
	httpClient := http.NewClientHttp(host, port)
	demo := http.NewDemo(httpClient)
	demo.Exec()
}

//func registerUser() {
//	relPath := endpoints.EndpointRegisterUser
//	fmt.Printf("register new user...\n")
//	fmt.Printf("\t- request to: %s\n", relPath)
//
//	cmd := register.ClientArgs{
//		Uuid:       uuid.New().String(),
//		Alias:      "demo_alias",
//		Name:       "demo_name",
//		SecondName: "demo_second_name",
//		Email:      "demo_email@gmail.com",
//		Password:   "Linux64bits$",
//	}
//	bytesCmd, err := json.Marshal(cmd)
//	if err != nil {
//		log.Panic(err)
//	}
//
//	fmt.Printf("\t- request command:\n")
//	fmt.Printf("\t\t- Uuid: %v\n", cmd.Uuid)
//	fmt.Printf("\t\t- Alias: %v\n", cmd.Alias)
//	fmt.Printf("\t\t- Name: %v\n", cmd.Name)
//	fmt.Printf("\t\t- SecondName: %v\n", cmd.SecondName)
//	fmt.Printf("\t\t- Email: %v\n", cmd.Email)
//	fmt.Printf("\t\t- Password: %v\n", cmd.Password)
//
//	doRequest(http.MethodPost, relPath, bytesCmd)
//}
//
//func checkHealth() {
//	relPath := endpoints.EndpointCheckStatus
//	fmt.Printf("cheking App Status...\n")
//	fmt.Printf("\t- request to: %s\n", relPath)
//	doRequest(http.MethodGet, relPath, nil)
//}
//
//func doRequest(method, relPath string, body []byte) {
//	fullPath := fmt.Sprintf("http://0.0.0.0:8080%v", relPath)
//	request, err := http.NewRequest(method, fullPath, bytes.NewReader(body))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	request.Header.Set("Content-type", "application/json")
//
//	client := &http.Client{
//		Timeout: 200 * time.Second,
//	}
//
//	response, err := client.Do(request)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	respBody, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("\t- status code: %d\n", response.StatusCode)
//
//	switch method {
//	case http.MethodPost:
//		var gotRespBody *api.CommandResponse
//		if err := json.Unmarshal(respBody, &gotRespBody); err != nil {
//			log.Panicln(err)
//		}
//		fmt.Printf("\t- response body:\n")
//		printMessage(gotRespBody.Message)
//	case http.MethodGet:
//		var gotRespBody *api.QueryResponse
//		if err := json.Unmarshal(respBody, &gotRespBody); err != nil {
//			log.Panicln(err)
//		}
//		fmt.Printf("\t- response body:\n")
//		printMessage(gotRespBody.Message)
//	}
//}
//
//func printMessage(message message.MessageData) {
//	fmt.Printf("\t\tmessage:\n")
//	fmt.Printf("\t\t- ObjectId: %v\n", message.ObjectId)
//	fmt.Printf("\t\t- MessageId: %v\n", message.MessageId)
//	fmt.Printf("\t\t- MessagePkg: %v\n", message.MessagePkg)
//	fmt.Printf("\t\t- Variables: %v\n", message.Variables)
//	fmt.Printf("\t\t- Text: %v\n", message.Text)
//	fmt.Printf("\t\t- Time: %v\n", message.Time)
//	fmt.Printf("\t\t- ClientError: %v\n", message.ClientErrorType)
//}
