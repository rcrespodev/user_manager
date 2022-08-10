package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/api/v1/endpoints"
	"github.com/rcrespodev/user_manager/api/v1/handlers/getUser"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/test/integration"
	"github.com/rcrespodev/user_manager/test/integration/v1/handlers/ginHandlers/utils"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"testing"
	"time"
)

type args struct {
	token          string
	userUuid       string
	userEmail      string
	userAlias      string
	userName       string
	userSecondName string
}
type want struct {
	response       *api.QueryResponse
	httpStatusCode int
}

func TestGetUserGinHandlerFunc(t *testing.T) {
	userRepository := kernel.Instance.UserRepository()
	require.NoError(t, getUserSetup(userRepository))

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "correct request - all querystring match",
			args: args{
				token:          "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NTk3MjY0NDcsImtleSI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwMCJ9.qKsVOntsUC-PUSMQ7aiHsHDg0KrYBxZj19RgLAFqVOsVGq9Cuq0bYxFVZwfBkfiC5E0YZZ7qv1SCBKzkznl1DigQN-5XZlXDIau44bX5gRoY53oTWif342qU1q9-ZiSuDUuIG3TigfY8Uhf-gt3am_Re4pxnJCSrH4XdLDviS8U5XesmwI46ctsFkAJlDAiygpVHsfJy3iqwjWNrY2qrbUG3LZxe9QlfTKLeOapbPBSqktCQxlF6QaTKoAmrHlbbrAKNSaBXu9733E02m9-at6lz-klc8fomXTg2cIoskt2Cp1xHJ-kydqXfgoH0Zar90igifVT9fomh2Hmmp_hvPKSiAT2GtImMV2qjqY4AGBUWJ7qoFkYkvGOVEZS0OGPrPeVz62DbbLHHj4Mth26BLuH_sDKvGqK6d2OKeqhtvcdAPkAx5Gj7VtaeYUB8yhkJ01IzB2bmlVnIs6EtS8FQZMudNu-QmVt6_FYbi5-nmkCWOieWvstmV_gOjx-apv6Ie0ZRlKp6qycngSiuArXejmBwtCqvX7-vBJcgTNAa0hX1UEUkXhKDfLmXMcVJWMD79VOPHSM9WlO6czJe_sLMYxsb5LwiwJJw2e1e67DcBKw14AAUXIzzrIIyeNSU8LGuFU4_Rkuz6H7Ka3M--XlSIOz_LvooTrW8W7-rD5EgEr4",
				userUuid:       "123e4567-e89b-12d3-a456-426614174000",
				userEmail:      "email_exists@gmail.com",
				userAlias:      "alias_exists",
				userName:       "name_exists",
				userSecondName: "second_name_exists",
			},
			want: want{
				response: &api.QueryResponse{
					Message: message.MessageData{},
					Data: map[string]interface{}{
						"Uuid":           "123e4567-e89b-12d3-a456-426614174000",
						"Alias":          "alias_exists",
						"Name":           "Name_exists",
						"SecondName":     "Second_name_exists",
						"Email":          "email_exists@gmail.com",
						"HashedPassword": "",
					},
				},
				httpStatusCode: 200,
			},
		},
		{
			name: "correct request - single querystring match",
			args: args{
				token:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NTk3MjY0NDcsImtleSI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwMCJ9.qKsVOntsUC-PUSMQ7aiHsHDg0KrYBxZj19RgLAFqVOsVGq9Cuq0bYxFVZwfBkfiC5E0YZZ7qv1SCBKzkznl1DigQN-5XZlXDIau44bX5gRoY53oTWif342qU1q9-ZiSuDUuIG3TigfY8Uhf-gt3am_Re4pxnJCSrH4XdLDviS8U5XesmwI46ctsFkAJlDAiygpVHsfJy3iqwjWNrY2qrbUG3LZxe9QlfTKLeOapbPBSqktCQxlF6QaTKoAmrHlbbrAKNSaBXu9733E02m9-at6lz-klc8fomXTg2cIoskt2Cp1xHJ-kydqXfgoH0Zar90igifVT9fomh2Hmmp_hvPKSiAT2GtImMV2qjqY4AGBUWJ7qoFkYkvGOVEZS0OGPrPeVz62DbbLHHj4Mth26BLuH_sDKvGqK6d2OKeqhtvcdAPkAx5Gj7VtaeYUB8yhkJ01IzB2bmlVnIs6EtS8FQZMudNu-QmVt6_FYbi5-nmkCWOieWvstmV_gOjx-apv6Ie0ZRlKp6qycngSiuArXejmBwtCqvX7-vBJcgTNAa0hX1UEUkXhKDfLmXMcVJWMD79VOPHSM9WlO6czJe_sLMYxsb5LwiwJJw2e1e67DcBKw14AAUXIzzrIIyeNSU8LGuFU4_Rkuz6H7Ka3M--XlSIOz_LvooTrW8W7-rD5EgEr4",
				userUuid: "123e4567-e89b-12d3-a456-426614174000",
			},
			want: want{
				response: &api.QueryResponse{
					Message: message.MessageData{},
					Data: map[string]interface{}{
						"Uuid":           "123e4567-e89b-12d3-a456-426614174000",
						"Alias":          "alias_exists",
						"Name":           "Name_exists",
						"SecondName":     "Second_name_exists",
						"Email":          "email_exists@gmail.com",
						"HashedPassword": "",
					},
				},
				httpStatusCode: 200,
			},
		},
		{
			name: "bad request - userEmail querystring not match",
			args: args{
				token:          "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NTk3MjY0NDcsImtleSI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwMCJ9.qKsVOntsUC-PUSMQ7aiHsHDg0KrYBxZj19RgLAFqVOsVGq9Cuq0bYxFVZwfBkfiC5E0YZZ7qv1SCBKzkznl1DigQN-5XZlXDIau44bX5gRoY53oTWif342qU1q9-ZiSuDUuIG3TigfY8Uhf-gt3am_Re4pxnJCSrH4XdLDviS8U5XesmwI46ctsFkAJlDAiygpVHsfJy3iqwjWNrY2qrbUG3LZxe9QlfTKLeOapbPBSqktCQxlF6QaTKoAmrHlbbrAKNSaBXu9733E02m9-at6lz-klc8fomXTg2cIoskt2Cp1xHJ-kydqXfgoH0Zar90igifVT9fomh2Hmmp_hvPKSiAT2GtImMV2qjqY4AGBUWJ7qoFkYkvGOVEZS0OGPrPeVz62DbbLHHj4Mth26BLuH_sDKvGqK6d2OKeqhtvcdAPkAx5Gj7VtaeYUB8yhkJ01IzB2bmlVnIs6EtS8FQZMudNu-QmVt6_FYbi5-nmkCWOieWvstmV_gOjx-apv6Ie0ZRlKp6qycngSiuArXejmBwtCqvX7-vBJcgTNAa0hX1UEUkXhKDfLmXMcVJWMD79VOPHSM9WlO6czJe_sLMYxsb5LwiwJJw2e1e67DcBKw14AAUXIzzrIIyeNSU8LGuFU4_Rkuz6H7Ka3M--XlSIOz_LvooTrW8W7-rD5EgEr4",
				userUuid:       "123e4567-e89b-12d3-a456-426614174000",
				userEmail:      "email_not_exists@gmail.com",
				userAlias:      "alias_exists",
				userName:       "name_exists",
				userSecondName: "second_name_exists",
			},
			want: want{
				response: &api.QueryResponse{
					Message: message.MessageData{
						ObjectId:        "123e4567-e89b-12d3-a456-426614174000",
						MessageId:       17,
						MessagePkg:      "user",
						Text:            "none of the input values correspond to a registered user",
						ClientErrorType: 1,
					},
					Data: nil,
				},
				httpStatusCode: 400,
			},
		},
		{
			name: "bad request - querystring match but token dont belong to target user",
			args: args{
				token:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NTk3MjY0NDcsImtleSI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwMCJ9.qKsVOntsUC-PUSMQ7aiHsHDg0KrYBxZj19RgLAFqVOsVGq9Cuq0bYxFVZwfBkfiC5E0YZZ7qv1SCBKzkznl1DigQN-5XZlXDIau44bX5gRoY53oTWif342qU1q9-ZiSuDUuIG3TigfY8Uhf-gt3am_Re4pxnJCSrH4XdLDviS8U5XesmwI46ctsFkAJlDAiygpVHsfJy3iqwjWNrY2qrbUG3LZxe9QlfTKLeOapbPBSqktCQxlF6QaTKoAmrHlbbrAKNSaBXu9733E02m9-at6lz-klc8fomXTg2cIoskt2Cp1xHJ-kydqXfgoH0Zar90igifVT9fomh2Hmmp_hvPKSiAT2GtImMV2qjqY4AGBUWJ7qoFkYkvGOVEZS0OGPrPeVz62DbbLHHj4Mth26BLuH_sDKvGqK6d2OKeqhtvcdAPkAx5Gj7VtaeYUB8yhkJ01IzB2bmlVnIs6EtS8FQZMudNu-QmVt6_FYbi5-nmkCWOieWvstmV_gOjx-apv6Ie0ZRlKp6qycngSiuArXejmBwtCqvX7-vBJcgTNAa0hX1UEUkXhKDfLmXMcVJWMD79VOPHSM9WlO6czJe_sLMYxsb5LwiwJJw2e1e67DcBKw14AAUXIzzrIIyeNSU8LGuFU4_Rkuz6H7Ka3M--XlSIOz_LvooTrW8W7-rD5EgEr4",
				userUuid: "123e4567-e89b-12d3-a456-426614174001",
			},
			want: want{
				response: &api.QueryResponse{
					Message: message.MessageData{
						ObjectId:        "123e4567-e89b-12d3-a456-426614174000",
						MessageId:       0,
						MessagePkg:      "Authorization",
						Text:            "Unauthorized",
						ClientErrorType: 2,
					},
					Data: nil,
				},
				httpStatusCode: 401,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := integration.NewTestServerHttpGin(endpoints.Endpoints{
				endpoints.EndpointGetUser: endpoints.Endpoint{
					HttpMethod:     http.MethodGet,
					Handler:        getUser.GetUserGinHandlerFunc(),
					AuthValidation: true,
				},
			})

			response := mockServer.DoRequest(integration.DoRequestCommand{
				BodyRequest:  nil,
				RelativePath: endpoints.EndpointGetUser,
				Uuid:         tt.args.userUuid,
				Token:        tt.args.token,
				QueryString:  buildQueryString(tt.args),
			})

			// Header check
			require.EqualValues(t, tt.want.httpStatusCode, response.HttpCode)

			// Body check
			var gotRespBody *api.QueryResponse
			if err := json.Unmarshal(response.Body, &gotRespBody); err != nil {
				log.Panicln(err)
			}
			gotRespBody.Message.Time = time.Time{}

			require.EqualValues(t, tt.want.response, gotRespBody)

			// Jwt Check
			if response.HttpCode == 200 {
				utils.TokenValidationForTesting(t, response.Header)
			}
		})
	}
}

func getUserSetup(repository domain.UserRepository) error {
	users := []*domain.NewUserCommand{
		{
			Uuid:       "123e4567-e89b-12d3-a456-426614174000",
			Alias:      "alias_exists",
			Name:       "name_exists",
			SecondName: "second_name_exists",
			Email:      "email_exists@gmail.com",
			Password:   "Linux64bits$",
			IgnorePass: false,
		},
		{
			Uuid:       "123e4567-e89b-12d3-a456-426614174001",
			Alias:      "alias_exists_2",
			Name:       "name_exists_2",
			SecondName: "second_name_exists_2",
			Email:      "email_exists@gmail.com_2",
			Password:   "Linux64bits$_2",
			IgnorePass: false,
		},
	}
	return utils.TableUsersSetup(users, userRepositoryInstance)
}

func buildQueryString(args args) string {
	var filters []struct {
		key   string
		value string
	}

	if args.userUuid != "" {
		filters = append(filters, struct {
			key   string
			value string
		}{key: "uuid", value: args.userUuid})
	}

	if args.userAlias != "" {
		filters = append(filters, struct {
			key   string
			value string
		}{key: "alias", value: args.userAlias})
	}

	if args.userEmail != "" {
		filters = append(filters, struct {
			key   string
			value string
		}{key: "email", value: args.userEmail})
	}

	if args.userName != "" {
		filters = append(filters, struct {
			key   string
			value string
		}{key: "name", value: args.userName})
	}

	if args.userSecondName != "" {
		filters = append(filters, struct {
			key   string
			value string
		}{key: "second_name", value: args.userSecondName})
	}

	if filters == nil {
		return ""
	}

	var queryString string
	for i, filter := range filters {
		s := fmt.Sprintf("%s=%s", filter.key, filter.value)
		if i == 0 {
			s = fmt.Sprintf("?%s", s)
		}

		if i != len(filters)-1 {
			s = fmt.Sprintf("%s&", s)
		}

		queryString = fmt.Sprintf("%s%s", queryString, s)
	}

	return queryString
}
