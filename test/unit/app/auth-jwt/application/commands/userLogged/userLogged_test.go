package userLogged

import (
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/pkg/app/authJwt/application/commands/userLogged"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/authJwt/domain"
	jwtRepository "github.com/rcrespodev/user_manager/pkg/app/authJwt/repository"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/valueObjects"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/repository"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var mockMessageRepository = repository.NewMockMessageRepository([]repository.MockData{
	{
		Id:              0,
		Pkg:             message.AuthorizationPkg,
		Text:            "Unauthorized",
		ClientErrorType: message.ClientErrorUnauthorized,
	},
})

var mockJwtRepository = jwtRepository.NewMockJwtRepository(jwtRepository.MockData{})

func TestUserLogged(t *testing.T) {
	type args struct {
		uuid           string
		expirationTime time.Duration
		publicKeyPath  string
		privateKeyPath string
	}
	type want struct {
		status         valueObjects.Status
		httpCodeReturn valueObjects.HttpCodeReturn
		error          error
		errorMessage   *message.MessageData
		successMessage *message.MessageData
		token          *jwtDomain.JwtSchema
	}
	test := []struct {
		name string
		args args
		want want
	}{
		{
			name: "correct keys",
			args: args{
				uuid:           "123e4567-e89b-12d3-a456-426614174000",
				expirationTime: 5,
				publicKeyPath:  "id_rsa.pub",
				privateKeyPath: "id_rsa",
			},
			want: want{
				status:         "",
				httpCodeReturn: 0,
				error:          nil,
				errorMessage:   nil,
				successMessage: nil,
				token: &jwtDomain.JwtSchema{
					Uuid:    "123e4567-e89b-12d3-a456-426614174000",
					IsValid: true,
					Token:   "",
				},
			},
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			cmdUuid, err := uuid.Parse(tt.args.uuid)
			require.NoError(t, err)

			retLog := domain.NewReturnLog(cmdUuid, mockMessageRepository, "")

			publicKey, err := os.ReadFile(tt.args.publicKeyPath)
			require.NoError(t, err)
			privateKey, err := os.ReadFile(tt.args.privateKeyPath)
			require.NoError(t, err)

			jwt := jwtDomain.NewJwt(publicKey, privateKey, tt.args.expirationTime)
			userLogger := userLogged.NewUserLogger(jwt, mockJwtRepository)

			userLoggedCmd := userLogged.NewCommand(cmdUuid)
			handler := userLogged.NewCommandHandler(userLogger)
			done := make(chan bool)
			go handler.Handle(userLoggedCmd, retLog, done)
			<-done

			// ReturnLog check
			require.EqualValues(t, tt.want.status, retLog.Status())
			require.EqualValues(t, tt.want.httpCodeReturn, retLog.HttpCode())

			// Check Internal error
			switch tt.want.error {
			case nil:
				if retLog.Error() != nil {
					require.Nil(t, retLog.Error().InternalError())
				}
			default:
				require.EqualValues(t, tt.want.error, retLog.Error().InternalError().Error())
			}

			// Check Client error messages
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

			// repository check
			token := mockJwtRepository.FindByUuid(jwtDomain.FindByUuidQuery{Uuid: cmdUuid.String()})
			switch tt.want.token {
			case nil:
				require.Nil(t, token)
			default:
				if tt.want.token.Token == "" {
					token.Token = ""
				}
				require.EqualValues(t, tt.want.token, token)
			}
		})
	}
}
