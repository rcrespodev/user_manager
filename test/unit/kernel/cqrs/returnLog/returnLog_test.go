package returnLog

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	domain "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/valueObjects"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/repository"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/service"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

const (
	defaultPkg = "testing"
)

func TestReturnLogSrv(t *testing.T) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	homeProject := os.Getenv("HOME_PROJECT")

	var mockRepository = repository.NewMockMessageRepository([]repository.MockData{
		{
			Id:   001,
			Pkg:  "testing",
			Text: "message 001, pkg testing with vars %v, %v.",
		},
		{
			Id:   002,
			Pkg:  "testing",
			Text: "message 002, pkg testing with vars %v, %v, %v, %v.",
		},
	})

	type internalError struct {
		Error error
		file  string
	}
	type wantCustomError struct {
		clienteError *message.MessageData
		*internalError
	}
	type wantSuccess struct {
		message *message.MessageData
	}
	type args struct {
		defaultPkg string
		error      *domain.NewErrorCommand
		success    domain.NewSuccessCommand
	}
	type want struct {
		status         valueObjects.Status
		error          *wantCustomError
		success        *wantSuccess
		httpCodeReturn valueObjects.HttpCodeReturn
	}
	tests := []struct {
		name string
		args *args
		want *want
	}{
		{
			name: "message not found",
			args: &args{
				defaultPkg: defaultPkg,
				error:      nil,
				success: &message.NewMessageCommand{
					MessageId: 999,
					ObjectId:  "reference",
					Variables: message.Variables{"var1", "var2"},
				},
			},
			want: &want{
				status: valueObjects.Error,
				error: &wantCustomError{
					internalError: &internalError{
						Error: fmt.Errorf("message 999 not found in pkg testing"),
						file:  fmt.Sprintf("%v/user_manager/test/unit/kernel/cqrs/returnLog/returnLog_test.go", homeProject),
					},
				},
				success:        &wantSuccess{message: nil},
				httpCodeReturn: 500,
			},
		},
		{
			name: "one success and one internal error",
			args: &args{
				defaultPkg: defaultPkg,
				error: &domain.NewErrorCommand{
					Error:             fmt.Errorf("internal Error"),
					NewMessageCommand: nil,
				},
				success: &message.NewMessageCommand{
					MessageId: 001,
					ObjectId:  "reference",
					Variables: message.Variables{"var1", "var2"},
				},
			},
			want: &want{
				status: valueObjects.Error,
				error: &wantCustomError{
					internalError: &internalError{
						Error: fmt.Errorf("internal Error"),
						file:  fmt.Sprintf("%v/user_manager/test/unit/kernel/cqrs/returnLog/returnLog_test.go", homeProject),
					},
				},
				success:        &wantSuccess{message: nil},
				httpCodeReturn: 500,
			},
		},
		{
			name: "one success and one external error",
			args: &args{
				defaultPkg: defaultPkg,
				error: &domain.NewErrorCommand{
					Error: nil,
					NewMessageCommand: &message.NewMessageCommand{
						ObjectId:  "reference",
						MessageId: 001,
						Variables: message.Variables{"var1", "var2"},
					},
				},
				success: &message.NewMessageCommand{
					MessageId: 001,
					ObjectId:  "reference",
					Variables: message.Variables{"var1", "var2"},
				},
			},
			want: &want{
				status: valueObjects.Error,
				error: &wantCustomError{
					internalError: nil,
					clienteError: &message.MessageData{
						ObjectId:   "reference",
						MessageId:  001,
						MessagePkg: "testing",
						Variables:  message.Variables{"var1", "var2"},
						Text:       "message 001, pkg testing with vars var1, var2.",
						Time:       time.Time{},
					},
				},
				success:        &wantSuccess{message: nil},
				httpCodeReturn: 400,
			},
		},
		{
			name: "one internal error and one external error",
			args: &args{
				defaultPkg: defaultPkg,
				error: &domain.NewErrorCommand{
					Error: fmt.Errorf("internal Error"),
					NewMessageCommand: &message.NewMessageCommand{
						ObjectId:  "reference",
						MessageId: 001,
						Variables: message.Variables{"var1", "var2"},
					},
				},
				success: nil,
			},
			want: &want{
				status: valueObjects.Error,
				error: &wantCustomError{
					internalError: &internalError{
						Error: fmt.Errorf("internal Error"),
						file:  fmt.Sprintf("%v/user_manager/test/unit/kernel/cqrs/returnLog/returnLog_test.go", homeProject),
					},
					clienteError: nil,
				},
				success:        &wantSuccess{message: nil},
				httpCodeReturn: 500,
			},
		},
		{
			name: "external error and success",
			args: &args{
				defaultPkg: defaultPkg,
				error: &domain.NewErrorCommand{
					Error: nil,
					NewMessageCommand: &message.NewMessageCommand{
						ObjectId:  "reference error",
						MessageId: 001,
						Variables: message.Variables{"var1", "var2"},
					},
				},
				success: &message.NewMessageCommand{
					ObjectId:  "reference success",
					MessageId: 001,
					Variables: message.Variables{"var1", "var2"},
				},
			},
			want: &want{
				status: valueObjects.Error,
				error: &wantCustomError{
					internalError: nil,
					clienteError: &message.MessageData{
						ObjectId:   "reference error",
						MessageId:  001,
						MessagePkg: "testing",
						Variables:  message.Variables{"var1", "var2"},
						Text:       "message 001, pkg testing with vars var1, var2.",
					},
				},
				success:        &wantSuccess{message: nil},
				httpCodeReturn: 400,
			},
		},
		{
			name: "success",
			args: &args{
				defaultPkg: defaultPkg,
				error:      nil,
				success: &message.NewMessageCommand{
					ObjectId:  "reference success",
					MessageId: 002,
					Variables: message.Variables{"var1", "var2", "var3", "var4"},
				},
			},
			want: &want{
				status: valueObjects.Success,
				error:  nil,
				success: &wantSuccess{message: &message.MessageData{
					ObjectId:   "reference success",
					MessageId:  002,
					MessagePkg: "testing",
					Variables:  message.Variables{"var1", "var2", "var3", "var4"},
					Text:       "message 002, pkg testing with vars var1, var2, var3, var4.",
					Time:       time.Time{},
				}},
				httpCodeReturn: 200,
			},
		},
		{
			name: "more variables than msg vars",
			args: &args{
				defaultPkg: defaultPkg,
				error:      nil,
				success: &message.NewMessageCommand{
					ObjectId:  "reference success",
					MessageId: 001,
					Variables: message.Variables{"var1", "var2", "var3", "var4"},
				},
			},
			want: &want{
				status: valueObjects.Success,
				error:  nil,
				success: &wantSuccess{message: &message.MessageData{
					ObjectId:   "reference success",
					MessageId:  001,
					MessagePkg: "testing",
					Variables:  message.Variables{"var1", "var2", "var3", "var4"},
					Text:       "message 001, pkg testing with vars var1, var2.",
					Time:       time.Time{},
				}},
				httpCodeReturn: 200,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var srv service.ReturnLogSrv
			srv = domain.NewReturnLog(uuid.New(), mockRepository, tt.args.defaultPkg)

			if tt.args.error != nil {
				srv.LogError(*tt.args.error)
			}

			if tt.args.success != nil {
				srv.LogSuccess(tt.args.success)
			}

			errorLog := srv.Error()
			switch errorLog {
			case nil:
				if tt.want.error != nil {
					t.Errorf("Error()\n\t- got: %v\n\t- want: %v", errorLog, tt.want.error)
				}
			default:
				gotErrorMsg := errorLog.Message()
				if gotErrorMsg != nil {
					gotErrorMsg.Time = time.Time{} // clear time field
				}
				if !reflect.DeepEqual(gotErrorMsg, tt.want.error.clienteError) {
					t.Errorf("ErrorMessage()\n\t- got: %v\n\t- want: %v", gotErrorMsg, tt.want.error.clienteError)
				}

				internalError := errorLog.InternalError()
				switch internalError {
				case nil:
					if tt.want.error.internalError != nil {
						t.Errorf("Internal Error\n\t- got: %v\n\t- want: %v", internalError, tt.want.error.internalError)
					}
				default:
					if gotInternalError := internalError.Error(); !reflect.DeepEqual(gotInternalError, tt.want.error.internalError.Error) {
						t.Errorf("Internal Error\n\t- got: %v\n\t- want: %v", gotInternalError, tt.want.error.internalError.Error)
					}
					if gotInternalErrorFile := internalError.File(); !reflect.DeepEqual(gotInternalErrorFile, tt.want.error.internalError.file) {
						t.Errorf("Internal Error File\n\t- got: %v\n\t- want: %v", gotInternalErrorFile, tt.want.error.internalError.file)
					}
				}
			}

			gotSuccess := srv.Success()
			switch gotSuccess {
			case nil:
				if tt.want.success.message != nil {
					t.Errorf("Success()\n\t- got: %v\n\t- want: %v", gotSuccess, tt.want.success.message)
				}
			default:
				gotSuccessMsg := gotSuccess.MessageData()
				if gotSuccessMsg != nil {
					gotSuccessMsg.Time = time.Time{}
				}
				if !reflect.DeepEqual(gotSuccessMsg, tt.want.success.message) {
					t.Errorf("Success()\n\t- got: %v\n\t- want: %v", gotSuccess, tt.want.success.message)
				}
			}

			if gotHttpCode := srv.HttpCode(); !reflect.DeepEqual(gotHttpCode, tt.want.httpCodeReturn) {
				t.Errorf("HttpCode()\n\t- got: %v\n\t- want: %v", gotHttpCode, tt.want.httpCodeReturn)
			}

			if gotStatus := srv.Status(); !reflect.DeepEqual(gotStatus, tt.want.status) {
				t.Errorf("Status()\n\t- got: %v\n\t- want: %v", gotStatus, tt.want.status)
			}

		})
	}
}
