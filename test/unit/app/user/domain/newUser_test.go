package service

import (
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/valueObjects"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/repository"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	var messageRepository = repository.NewMockMessageRepository([]repository.MockData{
		// 001 --> user created, 002 --> user updated, 003 --> user deleted
		{
			Id:              004,
			Pkg:             "user",
			Text:            "value %v is invalid as %v attribute",
			ClientErrorType: message.ClientErrorBadRequest,
		},
		{
			Id:              005,
			Pkg:             "user",
			Text:            "attribute %v are mandatory",
			ClientErrorType: message.ClientErrorBadRequest,
		},
		{
			Id:              006,
			Pkg:             "user",
			Text:            "attribute %v can´t be greater than %v characters",
			ClientErrorType: message.ClientErrorBadRequest,
		},
		{
			Id:              007,
			Pkg:             "user",
			Text:            "attribute %v can´t contain special characters (%v)",
			ClientErrorType: message.ClientErrorBadRequest,
		},
		{
			Id:              8,
			Pkg:             "user",
			Text:            "password must be contain at least one special character like %$#&",
			ClientErrorType: message.ClientErrorBadRequest,
		},
		{
			Id:              9,
			Pkg:             "user",
			Text:            "password must be contain at least one number",
			ClientErrorType: message.ClientErrorBadRequest,
		},
		{
			Id:              10,
			Pkg:             "user",
			Text:            "attribute %v can´t be smaller than %v characters",
			ClientErrorType: message.ClientErrorBadRequest,
		},
		{
			Id:              11,
			Pkg:             "user",
			Text:            "password must be contain at least one upper case",
			ClientErrorType: message.ClientErrorBadRequest,
		},
		{
			Id:              12,
			Pkg:             "user",
			Text:            "password must be contain at least one lower case",
			ClientErrorType: message.ClientErrorBadRequest,
		},
		{
			Id:              13,
			Pkg:             "user",
			Text:            "%v attribute dont´t must contain %v",
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
		// successes
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
			name: "good request - only trim spaces",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      " martin_fowler ",
				name:       " martin ",
				secondName: " fowler ",
				email:      " foo@test.com ",
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
		// nil fields
		{
			name: "bad request - missing uuid",
			args: &args{
				uuid:       "",
				alias:      "martin_fowler",
				name:       "martin",
				secondName: "fowler",
				email:      "foo@test.com",
				password:   "Linux648$",
			},
			want: &want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "martin_fowler",
					MessageId:       005,
					MessagePkg:      "user",
					Variables:       message.Variables{"Uuid"},
					Text:            "attribute Uuid are mandatory",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		{
			name: "bad request - missing alias",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "",
				name:       "martin",
				secondName: "fowler",
				email:      "foo@test.com",
				password:   "Linux648$",
			},
			want: &want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "",
					MessageId:       005,
					MessagePkg:      "user",
					Variables:       message.Variables{"Alias"},
					Text:            "attribute Alias are mandatory",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		{
			name: "bad request - missing name",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin_fowler",
				name:       "",
				secondName: "fowler",
				email:      "foo@test.com",
				password:   "Linux648$",
			},
			want: &want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "martin_fowler",
					MessageId:       005,
					MessagePkg:      "user",
					Variables:       message.Variables{"Name"},
					Text:            "attribute Name are mandatory",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		{
			name: "bad request - missing secondName",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin_fowler",
				name:       "martin",
				secondName: "",
				email:      "foo@test.com",
				password:   "Linux648$",
			},
			want: &want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "martin_fowler",
					MessageId:       005,
					MessagePkg:      "user",
					Variables:       message.Variables{"SecondName"},
					Text:            "attribute SecondName are mandatory",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		{
			name: "bad request - missing email",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin_fowler",
				name:       "martin",
				secondName: "fowler",
				email:      "",
				password:   "Linux648$",
			},
			want: &want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "martin_fowler",
					MessageId:       005,
					MessagePkg:      "user",
					Variables:       message.Variables{"Email"},
					Text:            "attribute Email are mandatory",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		{
			name: "bad request - missing password",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin_fowler",
				name:       "martin",
				secondName: "fowler",
				email:      "foo@test.com",
				password:   "",
			},
			want: &want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "martin_fowler",
					MessageId:       005,
					MessagePkg:      "user",
					Variables:       message.Variables{"Password"},
					Text:            "attribute Password are mandatory",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		// special characters validations
		{
			name: "bad request - Invalid alias (Special chars - sql injection)",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "OR 5=5",
				name:       "martin",
				secondName: "fowler",
				email:      "foo@test.com",
				password:   "Linux648$",
			},
			want: &want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "OR 5=5",
					MessageId:       007,
					MessagePkg:      "user",
					Variables:       message.Variables{"alias", "="},
					Text:            "attribute alias can´t contain special characters (=)",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		{
			name: "bad request - Invalid email (Special chars - sql injection)",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin_fowler",
				name:       "martin",
				secondName: "fowler",
				email:      "OR 5=5",
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
					Variables:       message.Variables{"OR 5=5", "email"},
					Text:            "value OR 5=5 is invalid as email attribute",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		{
			name: "bad request - Invalid name (Special chars)",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin_fowler",
				name:       "martin&",
				secondName: "fowler",
				email:      "foo@test.com",
				password:   "Linux648$",
			},
			want: &want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "martin_fowler",
					MessageId:       007,
					MessagePkg:      "user",
					Variables:       message.Variables{"name", "&"},
					Text:            "attribute name can´t contain special characters (&)",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		{
			name: "bad request - Invalid secondName (Special chars)",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin_fowler",
				name:       "martin",
				secondName: "fo%ler",
				email:      "foo@test.com",
				password:   "Linux648$",
			},
			want: &want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "martin_fowler",
					MessageId:       007,
					MessagePkg:      "user",
					Variables:       message.Variables{"name", "%"},
					Text:            "attribute name can´t contain special characters (%)",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		{
			name: "bad request - password not contain special an special character",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin_fowler",
				name:       "martin",
				secondName: "fowler",
				email:      "foo@test.com",
				password:   "Linux648",
			},
			want: &want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "martin_fowler",
					MessageId:       8,
					MessagePkg:      "user",
					Variables:       message.Variables{},
					Text:            "password must be contain at least one special character like %$#&",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		// number validations
		{
			name: "bad request - password not contain an number",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin_fowler",
				name:       "martin",
				secondName: "fowler",
				email:      "foo@test.com",
				password:   "Linux$",
			},
			want: &want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "martin_fowler",
					MessageId:       9,
					MessagePkg:      "user",
					Variables:       message.Variables{},
					Text:            "password must be contain at least one number",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		// len validations
		{
			name: "bad request - Invalid alias (Len > 30)",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin_fowlermartin_fowlermartin_fowler",
				name:       "martin",
				secondName: "fowler",
				email:      "foo@test.com",
				password:   "Linux648$",
			},
			want: &want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "martin_fowlermartin_fowlermartin_fowler",
					MessageId:       006,
					MessagePkg:      "user",
					Variables:       message.Variables{"alias", "30"},
					Text:            "attribute alias can´t be greater than 30 characters",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		{
			name: "bad request - Invalid name (Len > 50)",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin_fowler",
				name:       "martin martin martin martin martin martin martin martin martin",
				secondName: "fowler",
				email:      "foo@test.com",
				password:   "Linux648$",
			},
			want: &want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "martin_fowler",
					MessageId:       006,
					MessagePkg:      "user",
					Variables:       message.Variables{"name", "50"},
					Text:            "attribute name can´t be greater than 50 characters",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		{
			name: "bad request - Invalid second name (Len > 50)",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin_fowler",
				name:       "martin",
				secondName: "fowler fowler fowler fowler fowler fowler fowler fowler fowler",
				email:      "foo@test.com",
				password:   "Linux648$",
			},
			want: &want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "martin_fowler",
					MessageId:       006,
					MessagePkg:      "user",
					Variables:       message.Variables{"name", "50"},
					Text:            "attribute name can´t be greater than 50 characters",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		{
			name: "bad request - Invalid password (Len > 16)",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin_fowler",
				name:       "martin",
				secondName: "fowler",
				email:      "foo@test.com",
				password:   "Linux648$Linux648$",
			},
			want: &want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "martin_fowler",
					MessageId:       006,
					MessagePkg:      "user",
					Variables:       message.Variables{"password", "16"},
					Text:            "attribute password can´t be greater than 16 characters",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		{
			name: "bad request - Invalid password (Len < 8)",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin_fowler",
				name:       "martin",
				secondName: "fowler",
				email:      "foo@test.com",
				password:   "Linux6$",
			},
			want: &want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "martin_fowler",
					MessageId:       10,
					MessagePkg:      "user",
					Variables:       message.Variables{"password", "8"},
					Text:            "attribute password can´t be smaller than 8 characters",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		// invalid values
		{
			name: "bad request - invalid uuid",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-42661417",
				alias:      "martin_fowler",
				name:       "martin",
				secondName: "fowler",
				email:      "foo@test.com",
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
					Variables:       message.Variables{"123e4567-e89b-12d3-a456-42661417", "aggregateId"},
					Text:            "value 123e4567-e89b-12d3-a456-42661417 is invalid as aggregateId attribute",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		{
			name: "bad request - alias have an space",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin fowler",
				name:       "martin",
				secondName: "fowler",
				email:      "foo@test.com",
				password:   "Linux648$",
			},
			want: &want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "martin fowler",
					MessageId:       13,
					MessagePkg:      "user",
					Variables:       message.Variables{"alias", "spaces"},
					Text:            "alias attribute dont´t must contain spaces",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		{
			name: "bad request - Invalid email",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin_fowler",
				name:       "martin",
				secondName: "fowler",
				email:      "foo.test.com", //foo.test.com
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
					Variables:       message.Variables{"foo.test.com", "email"},
					Text:            "value foo.test.com is invalid as email attribute",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		{
			name: "bad request - password dont´t have upper case",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin_fowler",
				name:       "martin",
				secondName: "fowler",
				email:      "foo@test.com",
				password:   "linux648$",
			},
			want: &want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "martin_fowler",
					MessageId:       11,
					MessagePkg:      "user",
					Variables:       message.Variables{},
					Text:            "password must be contain at least one upper case",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
				userData:       nil,
			},
		},
		{
			name: "bad request - password dont´t have lower case",
			args: &args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin_fowler",
				name:       "martin",
				secondName: "fowler",
				email:      "foo@test.com",
				password:   "LINUX648$",
			},
			want: &want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "martin_fowler",
					MessageId:       12,
					MessagePkg:      "user",
					Variables:       message.Variables{},
					Text:            "password must be contain at least one lower case",
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
			byteUuid, _ := uuid.Parse(tt.args.uuid)
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
			switch tt.want.userData {
			case nil:
				require.Nil(t, user)
			default:
				require.EqualValues(t, tt.want.userData.Email, user.Email().Address())
				require.EqualValues(t, tt.want.userData.Name, user.Name().Name())
				require.EqualValues(t, tt.want.userData.SecondName, user.SecondName().Name())
				require.EqualValues(t, tt.want.userData.Uuid, user.Uuid().String())
				require.EqualValues(t, tt.want.userData.Alias, user.Alias().Alias())
				require.EqualValues(t, tt.want.userData.Password, user.Password().String())
				require.NotNil(t, string(user.Password().Hash()))
				err := bcrypt.CompareHashAndPassword(user.Password().Hash(), []byte(user.Password().String()))
				require.NoError(t, err)
			}

			// ReturnLog check
			require.EqualValues(t, retLog.Status(), tt.want.status)
			require.EqualValues(t, tt.want.httpCodeReturn, retLog.HttpCode())

			// Check Client error messages
			switch tt.want.error {
			case nil:
				if retLog.Error() != nil {
					require.Nil(t, retLog.Error().InternalError())
				}
			default:
				require.EqualValues(t, tt.want.error, retLog.Error().InternalError().Error())
			}

			switch tt.want.errorMessage {
			case nil:
				if retLog.Error() != nil {
					require.Nil(t, retLog.Error().Message())
				}
			default:
				gotMessage := retLog.Error().Message()
				gotMessage.Time = time.Time{}
				require.EqualValues(t, tt.want.errorMessage, gotMessage)
			}

			// Check Success message
			switch tt.want.successMessage {
			case nil:
				require.Nil(t, retLog.Success())
			default:
				gotMessage := retLog.Success().MessageData()
				gotMessage.Time = tt.want.successMessage.Time
				require.EqualValues(t, tt.want.successMessage, gotMessage)
			}
		})
	}
}
