package repository

import (
	"database/sql"
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/repository/mySql"
	"strings"
)

type MySqlUserRepository struct {
	mySqlRepository *mySql.MySqlRepository
	trx             *sql.Tx
}

func NewMySqlUserRepository(mySqlClient *sql.DB) *MySqlUserRepository {
	return &MySqlUserRepository{
		mySqlRepository: mySql.NewMySqlRepository(mySqlClient),
	}
}

func (m *MySqlUserRepository) SaveUser(user *domain.User, log *returnLog.ReturnLog) {
	const (
		saveUser = "INSERT INTO users (uuid, alias, name, second_name, email, password) values (?, ?, ?, ?, ?, ?);"
	)

	if err := m.newTrx(); err != nil {
		log.LogError(returnLog.NewErrorCommand{Error: err})
		return
	}

	resp, err := m.trx.Exec(saveUser, user.Uuid().String(), user.Alias().Alias(), user.Name().Name(),
		user.SecondName().Name(), user.Email().Address(), string(user.Password().Hash()))
	if err != nil {
		_ = m.trx.Rollback()
		log.LogError(returnLog.NewErrorCommand{Error: err})
		return
	}
	if rows, err := resp.RowsAffected(); rows == 0 || err != nil {
		log.LogError(returnLog.NewErrorCommand{Error: err})
	}

	if err := m.trx.Commit(); err != nil {
		_ = m.trx.Rollback()
		log.LogError(returnLog.NewErrorCommand{Error: err})
		return
	}

	return
}

func (m *MySqlUserRepository) FindUser(command domain.FindUserCommand) *domain.UserSchema {
	whereValues, whereFields := make([]interface{}, len(command.Where)), make([]string, len(command.Where))

	for i, args := range command.Where {
		whereFields[i] = fmt.Sprintf("%s = ?", args.Field)
		whereValues[i] = args.Value
	}
	queryString := fmt.Sprintf("SELECT uuid, alias, name, second_name, email, password FROM users WHERE %s;",
		strings.Join(whereFields, " AND "))

	userSchema := &domain.UserSchema{}
	if err := m.newTrx(); err != nil {
		command.Log.LogError(returnLog.NewErrorCommand{Error: err})
		return nil
	}

	err := m.trx.QueryRow(queryString, whereValues...).Scan(&userSchema.Uuid, &userSchema.Alias,
		&userSchema.Name, &userSchema.SecondName, &userSchema.Email, &userSchema.HashedPassword)
	if err == sql.ErrNoRows || userSchema.Uuid == "" {
		return nil
	}

	return userSchema
	//return domain.NewUser(domain.NewUserCommand{
	//	Uuid:       userSchema.Uuid,
	//	Alias:      userSchema.Alias,
	//	Name:       userSchema.Name,
	//	SecondName: userSchema.SecondName,
	//	Email:      userSchema.Email,
	//	Password:   command.Password,
	//}, command.Log)
}

func (m *MySqlUserRepository) newTrx() error {
	trx, err := m.mySqlRepository.NewTrx()
	if err != nil {
		return err
	}

	m.trx = trx
	return nil
}
