package service

import (
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/valueObjects"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/repository"
	"log"
	"reflect"
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	var messageRepository = repository.NewMockMessageRepository([]repository.MockData{
		// 001 --> user created, 002 --> user updated, 003 --> user deleted
		{
			Id:              004,
			Pkg:             "user",
			Text:            "value foo.test.com is invalid as email attribute",
			ClientErrorType: message.ClientErrorBadRequest,
		},
	})

	type args struct {
		uuid       string
		alias      string
		name       string
		secondName string
		email      string
		password   string
	}
	type UserData struct {
		Uuid       string
		Alias      string
		Name       string
		SecondName string
		Email      string
		Password   string
	}
	type want struct {
		status         valueObjects.Status
		httpCodeReturn valueObjects.HttpCodeReturn
		error          error
		errorMessage   *message.MessageData
		successMessage *message.MessageData
		userData       *UserData
	}
	tests := []struct {
		name string
		args *args
		want *want
	}{
		{
			name: "good request",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin_fowler",
				name:       "martin",
				secondName: "fowler",
				email:      "foo@test.com",
				password:   "Linux648$",
			},
			want: &want{
				status:         "",
				httpCodeReturn: 000,
				error:          nil,
				errorMessage:   nil,
				successMessage: nil,
				userData: &UserData{
					Uuid:       "123e4567-e89b-12d3-a456-426614174000",
					Alias:      "martin_fowler",
					Name:       "Martin",
					SecondName: "Fowler",
					Email:      "foo@test.com",
					Password:   "Linux648$",
				},
			},
		},
		{
			name: "bad request - Invalid email",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin_fowler",
				name:       "martin",
				secondName: "fowler",
				email:      "foo.test.com",
				password:   "Linux648$",
			},
			want: &want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "martin_fowler",
					MessageId:       004,
					MessagePkg:      "user",
					Variables:       message.Variables{"foo.test.com"},
					Text:            "value foo.test.com is invalid as email attribute",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			byteUuid, err := uuid.Parse(tt.args.uuid)
			if err != nil {
				log.Fatal(err)
			}
			retLog := returnLog.NewReturnLog(byteUuid, messageRepository, "user")
			user := domain.NewUser(domain.NewUserCommand{
				Uuid:       tt.args.uuid,
				Alias:      tt.args.alias,
				Name:       tt.args.name,
				SecondName: tt.args.secondName,
				Email:      tt.args.email,
				Password:   tt.args.password,
			}, retLog)

			// User Data check
			switch user {
			case nil:
				if tt.want.userData != nil {
					t.Errorf("User()\n\t- got: %v\n\t- want: %v", user, tt.want.userData)
				}
			default:
				if gotUserEmail := user.Email().Address(); !reflect.DeepEqual(gotUserEmail, tt.want.userData.Email) {
					t.Errorf("UserEmail()\n\t- got: %v\n\t- want: %v", gotUserEmail, tt.want.userData.Email)
				}
				if gotUserName := user.Name().Name(); !reflect.DeepEqual(gotUserName, tt.want.userData.Name) {
					t.Errorf("UserName()\n\t- got: %v\n\t- want: %v", gotUserName, tt.want.userData.Name)
				}
				if gotUserSecondName := user.SecondName().Name(); !reflect.DeepEqual(gotUserSecondName, tt.want.userData.SecondName) {
					t.Errorf("UserSecondName()\n\t- got: %v\n\t- want: %v", gotUserSecondName, tt.want.userData.SecondName)
				}
				if gotUserUuid := user.Uuid().String(); !reflect.DeepEqual(gotUserUuid, tt.want.userData.Uuid) {
					t.Errorf("UserUuid()\n\t- got: %v\n\t- want: %v", gotUserUuid, tt.want.userData.Uuid)
				}
				if gotUserAlias := user.Alias().Alias(); !reflect.DeepEqual(gotUserAlias, tt.want.userData.Alias) {
					t.Errorf("UserAlias()\n\t- got: %v\n\t- want: %v", gotUserAlias, tt.want.userData.Alias)
				}
				if gotUserPassword := user.Password().String(); !reflect.DeepEqual(gotUserPassword, tt.want.userData.Password) {
					t.Errorf("UserPasswordString()\n\t- got: %v\n\t- want: %v", gotUserPassword, tt.want.userData.Password)
				}
				if gotUserPasswordHash := user.Password().Hash(); string(gotUserPasswordHash) == "" {
					t.Errorf("UserPasswordHash is initial")
				}
			}

			// ReturnLog check
			if gotStatus := retLog.Status(); !reflect.DeepEqual(gotStatus, tt.want.status) {
				t.Errorf("Stauts()\n\t- got: %v\n\t- want: %v", gotStatus, tt.want.status)
			}

			if gotHttpCode := retLog.HttpCode(); !reflect.DeepEqual(gotHttpCode, tt.want.httpCodeReturn) {
				t.Errorf("HttpCode()\n\t- got: %v\n\t- want: %v", gotHttpCode, tt.want.httpCodeReturn)
			}

			if tt.want.error != nil {
				if gotError := retLog.Error().InternalError().Error(); !reflect.DeepEqual(gotError, tt.want.error) {
					t.Errorf("Error()\n\t- got: %v\n\t- want: %v", gotError, tt.want.error)
				}
			} else {
				if retLog.Error() != nil {
					if retLog.Error().InternalError() != nil {
						t.Errorf("Error()\n\t- got: %v\n\t- want: %v", retLog.Error().InternalError(), tt.want.error)
					}
				}
			}

			if tt.want.successMessage != nil {
				gotMessage := retLog.Success().MessageData()
				gotMessage.Time = time.Time{}
				if !reflect.DeepEqual(gotMessage, tt.want.successMessage) {
					t.Errorf("SuccessMessage()\n\t- got: %v\n\t- want: %v", gotMessage, tt.want.successMessage)
				}
			} else {
				if retLog.Success() != nil {
					t.Errorf("SuccessMessage()\n\t- got: %v\n\t- want: %v", retLog.Success(), tt.want.successMessage)
				}
			}

			if tt.want.errorMessage != nil {
				gotMessage := retLog.Error().Message()
				gotMessage.Time = time.Time{}
				if !reflect.DeepEqual(gotMessage, tt.want.errorMessage) {
					t.Errorf("ErrorMessage()\n\t- got: %v\n\t- want: %v", gotMessage, tt.want.errorMessage)
				}
			} else {
				if retLog.Error() != nil {
					if retLog.Error().Message() != nil {
						t.Errorf("ErrorMessage()\n\t- got: %v\n\t- want: %v", retLog.Error().Message(), tt.want.successMessage)
					}
				}
			}
		})
	}
}
