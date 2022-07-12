package repository

import (
	"database/sql"
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/repository/mySql"
	"log"
	"sync"
)

type MySqlUserRepository struct {
	mySqlRepository *mySql.MySqlRepository
	trx             *sql.Tx
	log             *returnLog.ReturnLog
	wg              *sync.WaitGroup
	err             error
}

type UserSchema struct {
	Uuid       string
	Alias      string
	Name       string
	SecondName string
	Email      string
	Password   string
}

func NewMySqlUserRepository(mySqlClient *sql.DB) *MySqlUserRepository {
	return &MySqlUserRepository{
		mySqlRepository: mySql.NewMySqlRepository(mySqlClient),
	}
}

func (m *MySqlUserRepository) SaveUser(user *domain.User, log *returnLog.ReturnLog, wg *sync.WaitGroup) {
	const saveUser = "INSERT INTO users (uuid, alias, name, second_name, email, password) values (?, ?, ?, ?, ?, ?);"

	if err := m.trxPrepare(log, wg); err != nil {
		m.errorHandler()
		return
	}
	result, err := m.trx.Exec(saveUser, user.Uuid().String(), user.Alias().Alias(), user.Name().Name(),
		user.SecondName().Name(), user.Email().Address(), string(user.Password().Hash()))
	if err != nil {
		m.err = err
		m.errorHandler()
		return
	}

	if rows, err := result.RowsAffected(); err != nil || rows == 0 {
		m.err = err
		if m.err == nil {
			m.err = fmt.Errorf("no rows afected")
		}
		m.errorHandler()
		return
	}
	_ = m.trx.Commit()
	//if err := m.trx.Commit(); err != nil {
	//	m.err = err
	//	m.errorHandler()
	//	return
	//}
}

func (m *MySqlUserRepository) FindUserById(command domain.FindByIdCommand) *domain.User {
	const findById = "SELECT * FROM users WHERE uuid = ?"

	if err := m.trxPrepare(command.FindUserCommand.Log, command.FindUserCommand.Wg); err != nil {
		m.errorHandler()
		return nil
	}
	return m.selectByComponent(findById, command.Uuid.String(), command.FindUserCommand.Password)
}

func (m *MySqlUserRepository) FindUserByEmail(command domain.FindByEmailCommand) *domain.User {
	const findByEmail = "SELECT * FROM users WHERE email = ?;"

	if err := m.trxPrepare(command.FindUserCommand.Log, command.FindUserCommand.Wg); err != nil {
		m.errorHandler()
		return nil
	}
	return m.selectByComponent(findByEmail, command.Email.Address(), command.FindUserCommand.Password)
}

func (m *MySqlUserRepository) FindUserByAlias(command domain.FindByAliasCommand) *domain.User {
	const findByAlias = "SELECT * FROM users WHERE alias = ?;"

	if err := m.trxPrepare(command.FindUserCommand.Log, command.FindUserCommand.Wg); err != nil {
		m.errorHandler()
		return nil
	}
	return m.selectByComponent(findByAlias, command.Alias.Alias(), command.FindUserCommand.Password)
}

func (m *MySqlUserRepository) selectByComponent(query string, arg any, password string) *domain.User {
	//defer func() {
	//	m.wg.Done()
	//}()
	user := &UserSchema{}
	err := m.trx.QueryRow(query, arg).Scan(&user.Uuid, &user.Alias, &user.Name, &user.SecondName, &user.Email, &user.Password)
	log.Printf("%v, %v, %v, %v, %v, %v", user.Uuid, user.Alias, user.Name, user.SecondName, user.Email, user.Password)
	if err != nil {
		//m.err = err
		//m.errorHandler()
		log.Printf("%v", err)
		m.wg.Done()
		return nil
	}

	if user.Uuid == "" {
		return nil
	}

	m.wg.Done()
	return domain.NewUser(domain.NewUserCommand{
		Uuid:       user.Uuid,
		Alias:      user.Alias,
		Name:       user.Name,
		SecondName: user.SecondName,
		Email:      user.Email,
		Password:   password,
	}, m.log)
}

func (m *MySqlUserRepository) trxPrepare(log *returnLog.ReturnLog, wg *sync.WaitGroup) error {
	m.wg = wg
	m.log = log
	m.newTrx()
	return m.err
}

func (m *MySqlUserRepository) newTrx() {
	trx, err := m.mySqlRepository.NewTrx()
	m.trx = trx
	m.err = err
}

func (m *MySqlUserRepository) errorHandler() {
	if m.trx != nil {
		_ = m.trx.Rollback()
		//if rollbackErr := m.trx.Rollback(); rollbackErr != nil {
		//	m.err = rollbackErr
		//}
	}

	m.log.LogError(returnLog.NewErrorCommand{
		Error:  m.err,
		Caller: 3,
	})
	m.wg.Done()
	return
}
