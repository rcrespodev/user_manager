package auth

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/pkg/app/auth/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	exampleJwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
)

func TestJwt(t *testing.T) {
	type args struct {
		uuid   uuid.UUID
		secret string
	}
	type want struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "valid token",
			args: args{
				uuid:   uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				secret: "my_secret",
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "invalid token",
			args: args{
				uuid:   uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				secret: "my_secret",
			},
			want: want{
				err: fmt.Errorf("can´t parse string. The token is invalid"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwtConfig := domain.NewJwtConfig(tt.args.secret)
			tokenString, err := domain.SignJwt(tt.args.uuid, jwtConfig)
			require.Nil(t, err)

			if tt.want.err != nil {
				tokenString = exampleJwt
			}
			err = domain.ParseJwt(tokenString, jwtConfig)
			require.EqualValues(t, tt.want.err, err)
		})
	}
}
