package tokenValidation

import (
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/pkg/app/auth-jwt/application/commands/tokenValidation"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/auth-jwt/domain"
	jwtRepository "github.com/rcrespodev/user_manager/pkg/app/auth-jwt/repository"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/valueObjects"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/repository"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var mockRepository = repository.NewMockMessageRepository([]repository.MockData{
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

func TestTokenValidation(t *testing.T) {
	type args struct {
		uuid           string
		publicKeyPath  string
		privateKeyPath string
		token          string
	}
	type want struct {
		status         valueObjects.Status
		httpCodeReturn valueObjects.HttpCodeReturn
		error          error
		errorMessage   *message.MessageData
		successMessage *message.MessageData
	}
	test := []struct {
		name string
		args args
		want want
	}{
		{
			name: "correct token",
			args: args{
				uuid:           "123e4567-e89b-12d3-a456-426614174000",
				token:          "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NTkzODQxOTAsImtleSI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwMCJ9.sqCpS8ISEnmvVCMQMNkWvmb7OAhM8HIuXE6L4rl2X2uvn5wXbfRD79FmjaST3nFekSeNtr9SwicVxU9VRnylKK5xexDHaoPoBZ4h7W6KYsYcJGfF-SN2dKAUGkwHUExFbkwGEgny9dFZKXMArBSLkzG8wWx5QSn38j3a4D1fvBKvLazTAEhP3SAfm2idlgnSGV3BvwATsIDlO_WCpToBgTre9pumKJXHGm43Tay2Zk6vYLyR1-2J2he9jkZec6pS4Dxe52PuVQnGiC6RFv6ttqEP42H9t_dTzfxcz_cwbj0fd5zURK_5F5zTI0aZfd7siACwleB7qAF9jgw6-vyh0kad25S0q5n2NGOK5YULfCbMGg3zDHOKhKEKuWu5aqKEVSUtTowZCl-LGsuU1a3ngKEBvEpCyS7lDRRA8zOCMIWTYbaJCYGsSWdb7d4f0N_jzHTj_8HCPwl8AxBtpCUYFFKEKthVRUz5sh3YGfTOVJVvXoFYUdU-jN4CpDrQiPbSEcFa8NzVdylch5WT06f_KMzJ-_yKW_DPG_6zYKJmZM5l7VtnPzgXRyVXQfOlsE0WWdigZv_2nE7WyhlRwHtrZknqLT9Dr5LA3nn1nN4EEAk90R2pKFfXOlIwbNaO4bufbEgRN6TkQB0eiu80-4didZtJEJMJVRpz1Qfn7fwpcsM",
				publicKeyPath:  "id_rsa.pub",
				privateKeyPath: "id_rsa",
			},
			want: want{
				status:         "",
				httpCodeReturn: 0,
				error:          nil,
				errorMessage:   nil,
				successMessage: nil,
			},
		},
		{
			name: "correct token but the token was invalidated",
			args: args{
				uuid:           "123e4567-e89b-12d3-a456-426614174001",
				token:          "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NTkzODcyNjIsImtleSI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwMSJ9.KGzLvI6GZg261KDLmWLpNf75d-7i_i1upvuvW6-_knqzRhvItcNR9nAyFO4AXM9Q2CQWnreaEqhrX5R49Q6vzIS19ykt5qpf7FgBrRVQQ5M7f05JkrqzvHBMe0yjGKPA3sdED1stLiwScBn8xhK7dEU7sj6WrYL7HxdNHyymZ5n7Neywpzp_bS_ssWYtAcxQq7waBR7V_eb633VgcOWVl0VpKQ4YfhqKN8i5YSUSASQUAF86MQ-f28dCuKYM4EVLenguDjl8-50mwBzwYKvFMRID0ys-EFlcm2jXvF3J6_QgoFeElGUIL7Hw4gCDwXrcmhRbJv1jN13O1zmQdL8cL5YVnoJBKLLggET2J8WmC5p3C0ILDQ2P60FcTtqLYV0uGOiGHthXvvq6MzXvcBebEGceV3XZvzouWTeowWvIrjZk-_ObfLxu5MctN7NFm1pKybjnwHVH_ToF3sf0pQoGEDQ-lRY0iOKIAlFBhYTpkFfs4VaNiRnR0GnVPog09JU9jFJN6YfCY5jjd21DOamS3hGHIts1vvJtRdGnUxCwFsnoaSgqnv76MMo9O9lm_G1CaBohdmz9uF-NsBf3N36zYNDuyqKxnaXI2_4fuLDnJ4qQ-6-y8gcFn9UQN3kX2nsekeqxgPHooaKgLNuvIUQRew-7oGp5acl9_eTiMMD38XY",
				publicKeyPath:  "id_rsa.pub",
				privateKeyPath: "id_rsa",
			},
			want: want{
				status:         valueObjects.Error,
				httpCodeReturn: 401,
				errorMessage: &message.MessageData{
					ObjectId:        "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NTkzODcyNjIsImtleSI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwMSJ9.KGzLvI6GZg261KDLmWLpNf75d-7i_i1upvuvW6-_knqzRhvItcNR9nAyFO4AXM9Q2CQWnreaEqhrX5R49Q6vzIS19ykt5qpf7FgBrRVQQ5M7f05JkrqzvHBMe0yjGKPA3sdED1stLiwScBn8xhK7dEU7sj6WrYL7HxdNHyymZ5n7Neywpzp_bS_ssWYtAcxQq7waBR7V_eb633VgcOWVl0VpKQ4YfhqKN8i5YSUSASQUAF86MQ-f28dCuKYM4EVLenguDjl8-50mwBzwYKvFMRID0ys-EFlcm2jXvF3J6_QgoFeElGUIL7Hw4gCDwXrcmhRbJv1jN13O1zmQdL8cL5YVnoJBKLLggET2J8WmC5p3C0ILDQ2P60FcTtqLYV0uGOiGHthXvvq6MzXvcBebEGceV3XZvzouWTeowWvIrjZk-_ObfLxu5MctN7NFm1pKybjnwHVH_ToF3sf0pQoGEDQ-lRY0iOKIAlFBhYTpkFfs4VaNiRnR0GnVPog09JU9jFJN6YfCY5jjd21DOamS3hGHIts1vvJtRdGnUxCwFsnoaSgqnv76MMo9O9lm_G1CaBohdmz9uF-NsBf3N36zYNDuyqKxnaXI2_4fuLDnJ4qQ-6-y8gcFn9UQN3kX2nsekeqxgPHooaKgLNuvIUQRew-7oGp5acl9_eTiMMD38XY",
					MessageId:       0,
					MessagePkg:      message.AuthorizationPkg,
					ClientErrorType: message.ClientErrorUnauthorized,
					Text:            "Unauthorized",
				},
			},
		},
		{
			name: "correct token but incorrect hash method",
			args: args{
				uuid:           "123e4567-e89b-12d3-a456-426614174000",
				token:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
				publicKeyPath:  "id_rsa.pub",
				privateKeyPath: "id_rsa",
			},
			want: want{
				status:         valueObjects.Error,
				httpCodeReturn: 401,
				errorMessage: &message.MessageData{
					ObjectId:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
					MessageId:       0,
					MessagePkg:      message.AuthorizationPkg,
					ClientErrorType: message.ClientErrorUnauthorized,
					Text:            "Unauthorized",
				},
			},
		},
		{
			name: "incorrect token",
			args: args{
				uuid:           "123e4567-e89b-12d3-a456-426614174000",
				token:          "eyJhbGciOiJSUzI1NiIsInR5cCI6Ikp.eyJleHAiOjE2NTkzODM5NTUsImlhdCI6MTY1OTM4MzY1NSwia2V5IjoiMTIzZTQ1NjctZTg5Yi0xMmQzLWE0NTYtNDI2NjE0MTc0MDAwIn0.nw9TSoAUCdqFstY60nXtCfYfKNFl95EO_I6GoTQ5S3waumtU5pzexy453qs3NhCbSQwFLuijso348Y26ptsxVzPYswTv356b6rGi9mgwg4hL1Sn6JXSTjdRLV7Z1rT1XS2YyNkS72ZOFBQlB2iJvRqtLNtu3EEWR6livEcBUHGR0zNkR8O5mVTQLR21Pi5Mmt1Ig0ONwzz5UkaeuAR3RrRQmavKQ2_ONbyR8C6F9Z6HYuIifZ4FuXLytvp4bUcpaQYyvjvGrgJOXlDcdQes8PZ8lnZtw2nyv4EtAMOnkmB4BaaXCD7jpm_L4aTpCHo-686uf_63bM6YseBKg8adLEPbQtbQUG_S4mZNw2ilUHKNg2sME3YsjPe3D50NhVCOIaUDNawp5zcH5VYSQfSjHI_PUhYGVu_ihyIOeKV7CaGj0pdKL9ByW8FgxJ9qbXYKiLM2KZIcwGCbmN6LLGQdKNfKWMhIv42enuPY8AcKStLFGrsq5QQ7dIOtkSMkJpWsL-id5ff5cD7m02mYKkhuVAYYqxgvJUlbRhvmo_ViwMioJHHmdHau3cE1BNiiugo2R3llqofv48Ul5dQudnROeik52l26ZEhR1OogHWVWAoarxQBon0ml7JtUZkebU7PY6hOCMFOho-at6I3coImKjedY0o2PwOSTdGteejzDdKDE",
				publicKeyPath:  "id_rsa.pub",
				privateKeyPath: "id_rsa",
			},
			want: want{
				status:         valueObjects.Error,
				httpCodeReturn: 401,
				errorMessage: &message.MessageData{
					ObjectId:        "eyJhbGciOiJSUzI1NiIsInR5cCI6Ikp.eyJleHAiOjE2NTkzODM5NTUsImlhdCI6MTY1OTM4MzY1NSwia2V5IjoiMTIzZTQ1NjctZTg5Yi0xMmQzLWE0NTYtNDI2NjE0MTc0MDAwIn0.nw9TSoAUCdqFstY60nXtCfYfKNFl95EO_I6GoTQ5S3waumtU5pzexy453qs3NhCbSQwFLuijso348Y26ptsxVzPYswTv356b6rGi9mgwg4hL1Sn6JXSTjdRLV7Z1rT1XS2YyNkS72ZOFBQlB2iJvRqtLNtu3EEWR6livEcBUHGR0zNkR8O5mVTQLR21Pi5Mmt1Ig0ONwzz5UkaeuAR3RrRQmavKQ2_ONbyR8C6F9Z6HYuIifZ4FuXLytvp4bUcpaQYyvjvGrgJOXlDcdQes8PZ8lnZtw2nyv4EtAMOnkmB4BaaXCD7jpm_L4aTpCHo-686uf_63bM6YseBKg8adLEPbQtbQUG_S4mZNw2ilUHKNg2sME3YsjPe3D50NhVCOIaUDNawp5zcH5VYSQfSjHI_PUhYGVu_ihyIOeKV7CaGj0pdKL9ByW8FgxJ9qbXYKiLM2KZIcwGCbmN6LLGQdKNfKWMhIv42enuPY8AcKStLFGrsq5QQ7dIOtkSMkJpWsL-id5ff5cD7m02mYKkhuVAYYqxgvJUlbRhvmo_ViwMioJHHmdHau3cE1BNiiugo2R3llqofv48Ul5dQudnROeik52l26ZEhR1OogHWVWAoarxQBon0ml7JtUZkebU7PY6hOCMFOho-at6I3coImKjedY0o2PwOSTdGteejzDdKDE",
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
			cmdUuid, err := uuid.Parse(tt.args.uuid)
			require.NoError(t, err)

			retLog := domain.NewReturnLog(cmdUuid, mockRepository, "")

			publicKey, err := os.ReadFile(tt.args.publicKeyPath)
			require.NoError(t, err)
			privateKey, err := os.ReadFile(tt.args.privateKeyPath)
			require.NoError(t, err)

			jwt := jwtDomain.NewJwt(publicKey, privateKey, 10)
			require.NoError(t, err)

			tokenValidator := tokenValidation.NewTokenValidator(jwt, mockJwtRepository)
			handler := tokenValidation.NewCommandHandler(tokenValidator)

			//token, err := jwt.CreateNewToken(cmdUuid.String())
			tokenValidationCmd := tokenValidation.NewCommand(tt.args.token)
			cmd := command.NewCommand(command.TokenValidation, cmdUuid, tokenValidationCmd)
			done := make(chan bool)
			go handler.Handle(*cmd, retLog, done)
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
		})
	}
}
