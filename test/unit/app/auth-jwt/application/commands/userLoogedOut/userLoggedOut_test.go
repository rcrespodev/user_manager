package userLoogedOut

import (
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/pkg/app/auth-jwt/application/commands/userLoggedOut"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/auth-jwt/domain"
	jwtRepository "github.com/rcrespodev/user_manager/pkg/app/auth-jwt/repository"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/valueObjects"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/repository"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var mockMessageRepository = repository.NewMockMessageRepository([]repository.MockData{
	{
		Id:              16,
		Pkg:             "user",
		Text:            "user logged out successful",
		ClientErrorType: 0,
	},
	{
		Id:              0,
		Pkg:             message.AuthorizationPkg,
		Text:            "Unauthorized",
		ClientErrorType: message.ClientErrorUnauthorized,
	},
})

var mockJwtRepository = jwtRepository.NewMockJwtRepository(jwtRepository.MockData{
	"123e4567-e89b-12d3-a456-426614174000": &jwtDomain.JwtSchema{
		Uuid:    "123e4567-e89b-12d3-a456-426614174000",
		IsValid: true,
		Token:   "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NTkzODQxOTAsImtleSI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwMCJ9.sqCpS8ISEnmvVCMQMNkWvmb7OAhM8HIuXE6L4rl2X2uvn5wXbfRD79FmjaST3nFekSeNtr9SwicVxU9VRnylKK5xexDHaoPoBZ4h7W6KYsYcJGfF-SN2dKAUGkwHUExFbkwGEgny9dFZKXMArBSLkzG8wWx5QSn38j3a4D1fvBKvLazTAEhP3SAfm2idlgnSGV3BvwATsIDlO_WCpToBgTre9pumKJXHGm43Tay2Zk6vYLyR1-2J2he9jkZec6pS4Dxe52PuVQnGiC6RFv6ttqEP42H9t_dTzfxcz_cwbj0fd5zURK_5F5zTI0aZfd7siACwleB7qAF9jgw6-vyh0kad25S0q5n2NGOK5YULfCbMGg3zDHOKhKEKuWu5aqKEVSUtTowZCl-LGsuU1a3ngKEBvEpCyS7lDRRA8zOCMIWTYbaJCYGsSWdb7d4f0N_jzHTj_8HCPwl8AxBtpCUYFFKEKthVRUz5sh3YGfTOVJVvXoFYUdU-jN4CpDrQiPbSEcFa8NzVdylch5WT06f_KMzJ-_yKW_DPG_6zYKJmZM5l7VtnPzgXRyVXQfOlsE0WWdigZv_2nE7WyhlRwHtrZknqLT9Dr5LA3nn1nN4EEAk90R2pKFfXOlIwbNaO4bufbEgRN6TkQB0eiu80-4didZtJEJMJVRpz1Qfn7fwpcsM",
	},
	"123e4567-e89b-12d3-a456-426614174001": &jwtDomain.JwtSchema{
		Uuid:    "123e4567-e89b-12d3-a456-426614174001",
		IsValid: false,
		Token:   "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NTkzODcyNjIsImtleSI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwMSJ9.KGzLvI6GZg261KDLmWLpNf75d-7i_i1upvuvW6-_knqzRhvItcNR9nAyFO4AXM9Q2CQWnreaEqhrX5R49Q6vzIS19ykt5qpf7FgBrRVQQ5M7f05JkrqzvHBMe0yjGKPA3sdED1stLiwScBn8xhK7dEU7sj6WrYL7HxdNHyymZ5n7Neywpzp_bS_ssWYtAcxQq7waBR7V_eb633VgcOWVl0VpKQ4YfhqKN8i5YSUSASQUAF86MQ-f28dCuKYM4EVLenguDjl8-50mwBzwYKvFMRID0ys-EFlcm2jXvF3J6_QgoFeElGUIL7Hw4gCDwXrcmhRbJv1jN13O1zmQdL8cL5YVnoJBKLLggET2J8WmC5p3C0ILDQ2P60FcTtqLYV0uGOiGHthXvvq6MzXvcBebEGceV3XZvzouWTeowWvIrjZk-_ObfLxu5MctN7NFm1pKybjnwHVH_ToF3sf0pQoGEDQ-lRY0iOKIAlFBhYTpkFfs4VaNiRnR0GnVPog09JU9jFJN6YfCY5jjd21DOamS3hGHIts1vvJtRdGnUxCwFsnoaSgqnv76MMo9O9lm_G1CaBohdmz9uF-NsBf3N36zYNDuyqKxnaXI2_4fuLDnJ4qQ-6-y8gcFn9UQN3kX2nsekeqxgPHooaKgLNuvIUQRew-7oGp5acl9_eTiMMD38XY",
	},
})

func TestUserLoggedOut(t *testing.T) {
	type args struct {
		uuid string
	}
	type want struct {
		status         valueObjects.Status
		httpCodeReturn valueObjects.HttpCodeReturn
		error          error
		errorMessage   *message.MessageData
		successMessage *message.MessageData
		jwtSchema      *jwtDomain.JwtSchema
	}
	test := []struct {
		name string
		args args
		want want
	}{
		{
			name: "logged out successful",
			args: args{
				uuid: "123e4567-e89b-12d3-a456-426614174000",
			},
			want: want{
				status:         valueObjects.Success,
				httpCodeReturn: 200,
				error:          nil,
				errorMessage:   nil,
				successMessage: &message.MessageData{
					ObjectId:        "123e4567-e89b-12d3-a456-426614174000",
					MessageId:       16,
					MessagePkg:      "user",
					Text:            "user logged out successful",
					ClientErrorType: 0,
				},
				jwtSchema: &jwtDomain.JwtSchema{
					Uuid:    "123e4567-e89b-12d3-a456-426614174000",
					IsValid: false,
					Token:   "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NTkzODQxOTAsImtleSI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwMCJ9.sqCpS8ISEnmvVCMQMNkWvmb7OAhM8HIuXE6L4rl2X2uvn5wXbfRD79FmjaST3nFekSeNtr9SwicVxU9VRnylKK5xexDHaoPoBZ4h7W6KYsYcJGfF-SN2dKAUGkwHUExFbkwGEgny9dFZKXMArBSLkzG8wWx5QSn38j3a4D1fvBKvLazTAEhP3SAfm2idlgnSGV3BvwATsIDlO_WCpToBgTre9pumKJXHGm43Tay2Zk6vYLyR1-2J2he9jkZec6pS4Dxe52PuVQnGiC6RFv6ttqEP42H9t_dTzfxcz_cwbj0fd5zURK_5F5zTI0aZfd7siACwleB7qAF9jgw6-vyh0kad25S0q5n2NGOK5YULfCbMGg3zDHOKhKEKuWu5aqKEVSUtTowZCl-LGsuU1a3ngKEBvEpCyS7lDRRA8zOCMIWTYbaJCYGsSWdb7d4f0N_jzHTj_8HCPwl8AxBtpCUYFFKEKthVRUz5sh3YGfTOVJVvXoFYUdU-jN4CpDrQiPbSEcFa8NzVdylch5WT06f_KMzJ-_yKW_DPG_6zYKJmZM5l7VtnPzgXRyVXQfOlsE0WWdigZv_2nE7WyhlRwHtrZknqLT9Dr5LA3nn1nN4EEAk90R2pKFfXOlIwbNaO4bufbEgRN6TkQB0eiu80-4didZtJEJMJVRpz1Qfn7fwpcsM",
				},
			},
		},
		{
			name: "logged out successful - token is already invalid",
			args: args{
				uuid: "123e4567-e89b-12d3-a456-426614174001",
			},
			want: want{
				status:         valueObjects.Success,
				httpCodeReturn: 200,
				error:          nil,
				errorMessage:   nil,
				successMessage: &message.MessageData{
					ObjectId:        "123e4567-e89b-12d3-a456-426614174001",
					MessageId:       16,
					MessagePkg:      "user",
					Text:            "user logged out successful",
					ClientErrorType: 0,
				},
				jwtSchema: &jwtDomain.JwtSchema{
					Uuid:    "123e4567-e89b-12d3-a456-426614174001",
					IsValid: false,
					Token:   "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NTkzODcyNjIsImtleSI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwMSJ9.KGzLvI6GZg261KDLmWLpNf75d-7i_i1upvuvW6-_knqzRhvItcNR9nAyFO4AXM9Q2CQWnreaEqhrX5R49Q6vzIS19ykt5qpf7FgBrRVQQ5M7f05JkrqzvHBMe0yjGKPA3sdED1stLiwScBn8xhK7dEU7sj6WrYL7HxdNHyymZ5n7Neywpzp_bS_ssWYtAcxQq7waBR7V_eb633VgcOWVl0VpKQ4YfhqKN8i5YSUSASQUAF86MQ-f28dCuKYM4EVLenguDjl8-50mwBzwYKvFMRID0ys-EFlcm2jXvF3J6_QgoFeElGUIL7Hw4gCDwXrcmhRbJv1jN13O1zmQdL8cL5YVnoJBKLLggET2J8WmC5p3C0ILDQ2P60FcTtqLYV0uGOiGHthXvvq6MzXvcBebEGceV3XZvzouWTeowWvIrjZk-_ObfLxu5MctN7NFm1pKybjnwHVH_ToF3sf0pQoGEDQ-lRY0iOKIAlFBhYTpkFfs4VaNiRnR0GnVPog09JU9jFJN6YfCY5jjd21DOamS3hGHIts1vvJtRdGnUxCwFsnoaSgqnv76MMo9O9lm_G1CaBohdmz9uF-NsBf3N36zYNDuyqKxnaXI2_4fuLDnJ4qQ-6-y8gcFn9UQN3kX2nsekeqxgPHooaKgLNuvIUQRew-7oGp5acl9_eTiMMD38XY",
				},
			},
		},
		{
			name: "token not found",
			args: args{
				uuid: "123e4567-e89b-12d3-a456-426614174002",
			},
			want: want{
				status:         valueObjects.Error,
				httpCodeReturn: 401,
				errorMessage: &message.MessageData{
					ObjectId:        "123e4567-e89b-12d3-a456-426614174002",
					MessageId:       0,
					MessagePkg:      message.AuthorizationPkg,
					ClientErrorType: message.ClientErrorUnauthorized,
					Text:            "Unauthorized",
				},
			},
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			uuidCmd, err := uuid.Parse(tt.args.uuid)
			require.NoError(t, err)

			retLog := domain.NewReturnLog(uuidCmd, mockMessageRepository, "user")

			userLoggedOutCmd := userLoggedOut.NewCommand(tt.args.uuid)
			cmd := command.NewCommand(command.UserLoggedOut, uuidCmd, userLoggedOutCmd)
			userLoggerOut := userLoggedOut.NewUserLoggerOut(mockJwtRepository)
			cmdHandler := userLoggedOut.NewCommandHandler(userLoggerOut)

			done := make(chan bool)
			go cmdHandler.Handle(*cmd, retLog, done)
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

			// Check repository
			targetJwt := mockJwtRepository.FindByUuid(jwtDomain.FindByUuidQuery{Uuid: tt.args.uuid})
			switch tt.want.jwtSchema {
			case nil:
				require.Nil(t, targetJwt)
			default:
				require.EqualValues(t, tt.want.jwtSchema, targetJwt)
			}
		})
	}
}
