package repository

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/repository"
	"sync"
)

type MySqlUserRepository struct {
	mySqlRepository *repository.MySqlRepository
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

func NewMySqlUserRepository() *MySqlUserRepository {
	return &MySqlUserRepository{
		mySqlRepository: repository.NewMySqlRepository(),
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
	if err := m.trx.Commit(); err != nil {
		m.err = err
		m.errorHandler()
		return
	}
}

func (m *MySqlUserRepository) FindUserById(uuid uuid.UUID, log *returnLog.ReturnLog, wg *sync.WaitGroup) *domain.User {
	const findById = "SELECT * FROM users WHERE uuid = ?"

	if err := m.trxPrepare(log, wg); err != nil {
		m.errorHandler()
		return nil
	}
	return m.selectByComponent(findById, uuid)
}

func (m *MySqlUserRepository) FindUserByEmail(email *domain.UserEmail, log *returnLog.ReturnLog, wg *sync.WaitGroup) *domain.User {
	const findByEmail = "SELECT * FROM users WHERE email = ?"

	if err := m.trxPrepare(log, wg); err != nil {
		m.errorHandler()
		return nil
	}
	return m.selectByComponent(findByEmail, email)
}

func (m *MySqlUserRepository) FindUserByAlias(alias *domain.UserAlias, log *returnLog.ReturnLog, wg *sync.WaitGroup) *domain.User {
	const findByAlias = "SELECT * FROM users WHERE alias = ?"

	if err := m.trxPrepare(log, wg); err != nil {
		m.errorHandler()
		return nil
	}
	return m.selectByComponent(findByAlias, alias)
}

func (m *MySqlUserRepository) selectByComponent(query string, arg any) *domain.User {
	var user UserSchema
	if err := m.trx.QueryRow(query, arg).Scan(&user); err != nil {
		m.err = err
		m.errorHandler()
	}

	return domain.NewUser(domain.NewUserCommand{
		Uuid:       user.Uuid,
		Alias:      user.Alias,
		Name:       user.Name,
		SecondName: user.SecondName,
		Email:      user.Email,
		Password:   user.Password,
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
		if rollbackErr := m.trx.Rollback(); rollbackErr != nil {
			m.err = rollbackErr
		}
	}

	m.log.LogError(returnLog.NewErrorCommand{
		Error: m.err,
	})
	m.wg.Done()
	return
}
