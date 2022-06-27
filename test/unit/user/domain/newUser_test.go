package service

import (
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/repository"
	"reflect"
	"testing"
)

func TestNewUser(t *testing.T) {
	var messageRepository = repository.NewMockMessageRepository([]repository.MockData{})

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
		userData *UserData
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
				userData: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retLog := returnLog.NewReturnLog(uuid.New(), messageRepository, "user")
			user := domain.NewUser(domain.NewUserCommand{
				Uuid:       tt.args.uuid,
				Alias:      tt.args.alias,
				Name:       tt.args.name,
				SecondName: tt.args.secondName,
				Email:      tt.args.email,
				Password:   tt.args.password,
			}, retLog)

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
		})
	}
}
