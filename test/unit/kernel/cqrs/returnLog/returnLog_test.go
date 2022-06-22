package returnLog

import (
	"fmt"
	"github.com/google/uuid"
	domain "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/valueObjects"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/repository"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/service"
	"os"
	"reflect"
	"testing"
)

func TestReturnLogSrv(t *testing.T) {
	homeProject := os.Getenv("HOME_PROJECT")
	const (
		defaultPkg = "testing"
	)
	var uuids [1]uuid.UUID
	for i, _ := range uuids {
		uuids[i] = uuid.New()
	}

	var mockRepository = repository.MockMessageRepository{}

	type internalError struct {
		Error error
		file  string
	}
	type wantCustomError struct {
		message *message.MessageData
		*internalError
	}
	type wantSuccess struct {
		message *message.MessageData
	}
	type args struct {
		uuid       uuid.UUID
		defaultPkg string
		error      *domain.NewErrorCommand
		success    *domain.NewSuccessCommand
	}
	type want struct {
		uuid           uuid.UUID
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
				uuid:       uuids[0],
				defaultPkg: defaultPkg,
				error:      nil,
				success: &domain.NewSuccessCommand{
					MessageId: 999,
					ObjectId:  "reference",
					Variables: message.Variables{"var1", "var2"},
				},
			},
			want: &want{
				uuid:   uuids[0],
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var srv service.ReturnLogSrv
			srv = domain.NewReturnLog(tt.args.uuid, mockRepository, tt.args.defaultPkg)

			if tt.args.error != nil {
				srv.LogError(domain.NewErrorCommand{
					Error: tt.args.error.Error,
					NewMessageCommand: message.NewMessageCommand{
						ObjectId:   tt.args.error.ObjectId,
						MessageId:  tt.args.error.MessageId,
						MessagePkg: tt.args.error.MessagePkg,
						Variables:  tt.args.error.Variables,
					},
				})
			}

			if tt.args.success != nil {
				srv.LogSuccess(domain.NewSuccessCommand{
					ObjectId:   tt.args.success.ObjectId,
					MessageId:  tt.args.success.MessageId,
					MessagePkg: tt.args.success.MessagePkg,
					Variables:  tt.args.success.Variables,
				})
			}

			errorLog := srv.Error()
			if gotErrorMsg := errorLog.Message(); !reflect.DeepEqual(gotErrorMsg, tt.want.error.message) {
				t.Errorf("ErrorMessage()\n\t- got: %v\n\t- want: %v", gotErrorMsg, tt.want.error.message)
			}

			internalError := errorLog.InternalError()
			if gotInternalError := internalError.Error(); !reflect.DeepEqual(gotInternalError, tt.want.error.internalError.Error) {
				t.Errorf("Internal Error\n\t- got: %v\n\t- want: %v", gotInternalError, tt.want.error.internalError.Error)
			}

			if gotInternalErrorFile := internalError.File(); !reflect.DeepEqual(gotInternalErrorFile, tt.want.error.internalError.file) {
				t.Errorf("Internal Error File\n\t- got: %v\n\t- want: %v", gotInternalErrorFile, tt.want.error.internalError.file)
			}

			gotSuccess := srv.Success()
			switch gotSuccess {
			case nil:
				if tt.want.success.message != nil {
					t.Errorf("Success()\n\t- got: %v\n\t- want: %v", gotSuccess, tt.want.success.message)
				}
			default:
				if !reflect.DeepEqual(gotSuccess.MessageData(), tt.want.success.message) {
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
