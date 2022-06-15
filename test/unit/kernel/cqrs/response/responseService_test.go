package response

import (
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/response/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/response/service"
	"reflect"
	"testing"
)

func TestResponseServicesType(t *testing.T) {
	type messageArgs struct {
		key        service.Key
		messagePkg string
		messageId  domain.MessageId
		variables  domain.Variables
		status     domain.MessageStatus
	}
	type wantLog struct {
		key            service.Key
		httpStatusCode domain.HttpStatusCode
		messages       domain.MessageBuffer
	}
	tests := []struct {
		name         string
		messages     []messageArgs
		wantLog      wantLog
		wantTotalLog []*wantLog
	}{
		{
			name: "Base test",
			messages: []messageArgs{
				{
					key:        "Base test",
					messagePkg: "testing",
					messageId:  001,
					variables:  domain.Variables{"var1", "var2", "var3", "var4"},
					status:     domain.Success,
				},
			},
			wantLog: wantLog{
				key:            "Base test",
				httpStatusCode: 200,
				messages: domain.MessageBuffer{
					{
						MessageId: 001,
						MessageData: domain.MessageData{
							MessageStatus: domain.Success,
							ErrorType:     0,
							Text:          "Success message with vars: var1, var2, var3.",
							Variables:     domain.Variables{"var1", "var2", "var3", "var4"},
						},
					},
				},
			},
			wantTotalLog: []*wantLog{
				{
					key:            "Base test",
					httpStatusCode: 200,
					messages: domain.MessageBuffer{
						{
							MessageId: 001,
							MessageData: domain.MessageData{
								MessageStatus: domain.Success,
								ErrorType:     0,
								Text:          "Success message with vars: var1, var2, var3.",
								Variables:     domain.Variables{"var1", "var2", "var3", "var4"},
							},
						},
					},
				},
			},
		},
	}

	var messageRepository *domain.MessageRepository
	srv := service.NewResponseService(messageRepository)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wantTotalLog []*wantLog
			for i, v := range tt.messages {
				switch i {
				case 0:
					srv.NewLog(v.key)
				default:
					if v.key != tt.messages[i-1].key {
						srv.NewLog(v.key)
					}
				}

				switch v.status {
				case domain.Error:
					srv.AddError(v.messageId, v.variables)
				case domain.Warning:
					srv.AddWarning(v.messageId, v.variables)
				case domain.Success:
					srv.AddSuccess(v.messageId, v.variables)
				}

				nextKey := v.key
				if len(tt.messages) > i+2 {
					nextKey = tt.messages[i+1].key
				}
				if v.key != nextKey || len(tt.messages) == i+1 {
					response := srv.GetCurrentResponse()

					httpCode := response.HttpStatusCode()
					if !reflect.DeepEqual(httpCode, tt.wantLog.httpStatusCode) {
						t.Errorf("Error in HttpStatusCode.\n\t- actual: %v \n\t- want:%v", httpCode, tt.wantLog.httpStatusCode)
						return
					}

					messages := response.MessageBuffer()
					for i, message := range messages {
						if !reflect.DeepEqual(message, tt.wantLog.messages[i]) {
							t.Errorf("Error in Messages.\n\t- actual: %v \n\t- want:%v", messages, tt.wantLog.messages[i])
							return
						}
					}

					appendLog := &wantLog{
						key:            v.key,
						httpStatusCode: httpCode,
						messages:       messages,
					}
					wantTotalLog = append(wantTotalLog, appendLog)
				}
			}
			if !reflect.DeepEqual(wantTotalLog, tt.wantTotalLog) {
				t.Errorf("Error in TotalLog.\n\t- actual: %v \n\t- want:%v", wantTotalLog, tt.wantTotalLog)
				return
			}
		})
	}
}
