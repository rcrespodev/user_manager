package handlers

import (
	"encoding/json"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/api/v1/endpoints"
	"github.com/rcrespodev/user_manager/api/v1/handlers/checkStatus"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/test/integration"
	"net/http"
	"reflect"
	"testing"
)

func TestCheckStatusGinHandlerFunc(t *testing.T) {
	type args struct {
		path string
	}
	type want struct {
		httpCode      int
		queryResponse api.QueryResponse
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "test",
			args: args{path: ""},
			want: want{
				httpCode: 200,
				queryResponse: api.QueryResponse{
					Message: message.MessageData{},
					Data:    "Check-Status = Ok",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := integration.NewTestServerHttpGin(endpoints.Endpoints{
				endpoints.EndpointCheckStatus: endpoints.Endpoint{
					HttpMethod: http.MethodGet,
					Handler:    checkStatus.StatusGinHandlerFunc(),
				},
			})
			response := testServer.DoRequest(integration.DoRequestCommand{
				BodyRequest:  nil,
				RelativePath: endpoints.EndpointCheckStatus,
			})
			var queryResponse api.QueryResponse
			err := json.Unmarshal(response.Body, &queryResponse)
			if err != nil {
				t.Fatal(err)
			}
			if gotHttpCode := response.HttpCode; !reflect.DeepEqual(gotHttpCode, tt.want.httpCode) {
				t.Errorf("Http Code\n\t- got: %v\n\t- want: %v", gotHttpCode, tt.want.httpCode)
			}
			if gotResponse := queryResponse; !reflect.DeepEqual(gotResponse, tt.want.queryResponse) {
				t.Errorf("Body Reponse\n\t- got: %v\n\t- want: %v", gotResponse, tt.want.queryResponse)
			}
		})
	}
}
