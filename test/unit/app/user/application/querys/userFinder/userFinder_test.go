package userFinder

import (
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/pkg/app/user/application/querys/userFinder"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/app/user/repository/userRepository"
	domain2 "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/repository"
	"github.com/stretchr/testify/require"
	"testing"
)

var messageRepository = repository.NewMockMessageRepository(nil)

func TestUserFinder(t *testing.T) {
	mockRepository := userRepository.NewMockUserRepository()
	setMockUserData(mockRepository)

	type args struct {
		query []domain.FindUserQuery
	}
	type want struct {
		userSchema *domain.UserSchema
	}
	test := []struct {
		name string
		args args
		want want
	}{
		{
			name: "find by uuid",
			args: args{
				query: []domain.FindUserQuery{
					{
						Log: domain2.NewReturnLog(uuid.New(), messageRepository, ""),
						Where: []domain.WhereArgs{
							{
								Field: "uuid",
								Value: "123e4567-e89b-12d3-a456-426614174000",
							},
						},
					},
				},
			},
			want: want{
				userSchema: &domain.UserSchema{
					Uuid:           "123e4567-e89b-12d3-a456-426614174000",
					Alias:          "user_exists",
					Name:           "Martin",
					SecondName:     "Fowler",
					Email:          "email_exists@test.com",
					HashedPassword: nil,
				},
			},
		},
		{
			name: "find by email",
			args: args{
				query: []domain.FindUserQuery{
					{
						Log: domain2.NewReturnLog(uuid.New(), messageRepository, ""),
						Where: []domain.WhereArgs{
							{
								Field: "email",
								Value: "email_exists@test.com",
							},
						},
					},
				},
			},
			want: want{
				userSchema: &domain.UserSchema{
					Uuid:           "123e4567-e89b-12d3-a456-426614174000",
					Alias:          "user_exists",
					Name:           "Martin",
					SecondName:     "Fowler",
					Email:          "email_exists@test.com",
					HashedPassword: nil,
				},
			},
		},
		{
			name: "find by alias",
			args: args{
				query: []domain.FindUserQuery{
					{
						Log: domain2.NewReturnLog(uuid.New(), messageRepository, ""),
						Where: []domain.WhereArgs{
							{
								Field: "alias",
								Value: "user_exists",
							},
						},
					},
				},
			},
			want: want{
				userSchema: &domain.UserSchema{
					Uuid:           "123e4567-e89b-12d3-a456-426614174000",
					Alias:          "user_exists",
					Name:           "Martin",
					SecondName:     "Fowler",
					Email:          "email_exists@test.com",
					HashedPassword: nil,
				},
			},
		},
		{
			name: "find by name",
			args: args{
				query: []domain.FindUserQuery{
					{
						Log: domain2.NewReturnLog(uuid.New(), messageRepository, ""),
						Where: []domain.WhereArgs{
							{
								Field: "name",
								Value: "Martin",
							},
						},
					},
				},
			},
			want: want{
				userSchema: &domain.UserSchema{
					Uuid:           "123e4567-e89b-12d3-a456-426614174000",
					Alias:          "user_exists",
					Name:           "Martin",
					SecondName:     "Fowler",
					Email:          "email_exists@test.com",
					HashedPassword: nil,
				},
			},
		},
		{
			name: "find by second name",
			args: args{
				query: []domain.FindUserQuery{
					{
						Log: domain2.NewReturnLog(uuid.New(), messageRepository, ""),
						Where: []domain.WhereArgs{
							{
								Field: "second_name",
								Value: "Fowler",
							},
						},
					},
				},
			},
			want: want{
				userSchema: &domain.UserSchema{
					Uuid:           "123e4567-e89b-12d3-a456-426614174000",
					Alias:          "user_exists",
					Name:           "Martin",
					SecondName:     "Fowler",
					Email:          "email_exists@test.com",
					HashedPassword: nil,
				},
			},
		},
		{
			name: "one query match",
			args: args{
				query: []domain.FindUserQuery{
					{
						Log: domain2.NewReturnLog(uuid.New(), messageRepository, ""),
						Where: []domain.WhereArgs{
							{
								Field: "email",
								Value: "email_exists@test.com",
							},
						},
					},
					{
						Log: domain2.NewReturnLog(uuid.New(), messageRepository, ""),
						Where: []domain.WhereArgs{
							{
								Field: "alias",
								Value: "alias_not_exists",
							},
						},
					},
					{
						Log: domain2.NewReturnLog(uuid.New(), messageRepository, ""),
						Where: []domain.WhereArgs{
							{
								Field: "second_name",
								Value: "second_name_not_exists",
							},
						},
					},
				},
			},
			want: want{
				userSchema: &domain.UserSchema{
					Uuid:           "123e4567-e89b-12d3-a456-426614174000",
					Alias:          "user_exists",
					Name:           "Martin",
					SecondName:     "Fowler",
					Email:          "email_exists@test.com",
					HashedPassword: nil,
				},
			},
		},
		{
			name: "two query match",
			args: args{
				query: []domain.FindUserQuery{
					{
						Log: domain2.NewReturnLog(uuid.New(), messageRepository, ""),
						Where: []domain.WhereArgs{
							{
								Field: "email",
								Value: "email_exists@test.com",
							},
						},
					},
					{
						Log: domain2.NewReturnLog(uuid.New(), messageRepository, ""),
						Where: []domain.WhereArgs{
							{
								Field: "alias",
								Value: "user_exists",
							},
						},
					},
					{
						Log: domain2.NewReturnLog(uuid.New(), messageRepository, ""),
						Where: []domain.WhereArgs{
							{
								Field: "second_name",
								Value: "second_name_not_exists",
							},
						},
					},
				},
			},
			want: want{
				userSchema: &domain.UserSchema{
					Uuid:           "123e4567-e89b-12d3-a456-426614174000",
					Alias:          "user_exists",
					Name:           "Martin",
					SecondName:     "Fowler",
					Email:          "email_exists@test.com",
					HashedPassword: nil,
				},
			},
		},
		{
			name: "not query match",
			args: args{
				query: []domain.FindUserQuery{
					{
						Log: domain2.NewReturnLog(uuid.New(), messageRepository, ""),
						Where: []domain.WhereArgs{
							{
								Field: "email",
								Value: "email_not_exists@test.com",
							},
						},
					},
					{
						Log: domain2.NewReturnLog(uuid.New(), messageRepository, ""),
						Where: []domain.WhereArgs{
							{
								Field: "alias",
								Value: "alias_not_exists",
							},
						},
					},
					{
						Log: domain2.NewReturnLog(uuid.New(), messageRepository, ""),
						Where: []domain.WhereArgs{
							{
								Field: "second_name",
								Value: "second_name_not_exists",
							},
						},
					},
				},
			},
			want: want{
				userSchema: nil,
			},
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			finder := userFinder.NewUserFinder(mockRepository)
			userSchema := finder.Exec(tt.args.query, domain2.NewReturnLog(uuid.New(), messageRepository, ""))
			if userSchema != nil {
				userSchema.HashedPassword = nil
			}
			require.EqualValues(t, tt.want.userSchema, userSchema)
		})
	}
}

func setMockUserData(userRepository *userRepository.MockUserRepository) {
	mockDataArgs := []domain.NewUserCommand{
		{
			Uuid:       "123e4567-e89b-12d3-a456-426614174000",
			Alias:      "user_exists",
			Name:       "martin",
			SecondName: "fowler",
			Email:      "email_exists@test.com",
			Password:   "Linux648$",
		},
	}
	var users []*domain.User
	for _, arg := range mockDataArgs {
		user := domain.NewUser(arg, domain2.NewReturnLog(uuid.New(), messageRepository, "user"))
		users = append(users, user)
	}
	userRepository.SetMockData(users)
}
