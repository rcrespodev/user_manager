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
	//log             *returnLog.ReturnLog
	//wg              *sync.WaitGroup
	//err  error
	//user chan domain.User
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

func (m *MySqlUserRepository) SaveUser(user *domain.User, log *returnLog.ReturnLog) {
	const saveUser = "INSERT INTO users (uuid, alias, name, second_name, email, password) values (?, ?, ?, ?, ?, ?);"

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

func (m *MySqlUserRepository) FindUser(command domain.FindUserCommand) *domain.User {
	whereValues, whereFields := make([]interface{}, len(command.Where)), make([]string, len(command.Where))

	for i, args := range command.Where {
		whereFields[i] = fmt.Sprintf("%s = ?", args.Field)
		whereValues[i] = args.Value
		//whereFields := append(whereFields, fmt.Sprintf("%s = ?", args.Field))
		//whereValues := append(whereValues, args.Value)
	}
	queryString := fmt.Sprintf("SELECT * FROM users WHERE %s;",
		strings.Join(whereFields, " AND "))

	userSchema := &UserSchema{}
	_ = m.newTrx()
	//rows, _ := m.trx.Query("SELECT * FROM users;")
	//if rows != nil {
	//	for rows.Next() {
	//		if err := rows.Scan(&userSchema.Uuid, &userSchema.Alias, &userSchema.Name,
	//			&userSchema.SecondName, &userSchema.Email, &userSchema.Password); err != nil {
	//			return nil
	//		}
	//	}
	//}

	err := m.trx.QueryRow(queryString, whereValues...).Scan(&userSchema.Uuid, &userSchema.Alias,
		&userSchema.Name, &userSchema.SecondName, &userSchema.Email, &userSchema.Password)
	if err == sql.ErrNoRows || userSchema.Uuid == "" {
		return nil
	}
	//if userSchema.Uuid == "" {
	//	return nil
	//}

	return domain.NewUser(domain.NewUserCommand{
		Uuid:       userSchema.Uuid,
		Alias:      userSchema.Alias,
		Name:       userSchema.Name,
		SecondName: userSchema.SecondName,
		Email:      userSchema.Email,
		Password:   command.Password,
	}, command.Log)
}

//func (m *MySqlUserRepository) FindUserById(command domain.FindByIdCommand, user chan *domain.User) {
//	const findById = "SELECT * FROM users WHERE uuid = ?"
//
//	if err := m.trxPrepare(command.FindUserCommand.Log); err != nil {
//		user <- nil
//		close(user)
//		m.errorHandler()
//		//m.wg.Done()
//		return
//	}
//	user <- m.selectByComponent(findById, command.Uuid.String(), command.FindUserCommand.Password)
//	//m.wg.Done()
//	close(user)
//	return
//}
//
//func (m *MySqlUserRepository) FindUserByEmail(command domain.FindByEmailCommand, user chan *domain.User) {
//	const findByEmail = "SELECT * FROM users WHERE email = ?;"
//
//	if err := m.trxPrepare(command.FindUserCommand.Log); err != nil {
//		user <- nil
//		close(user)
//		m.errorHandler()
//		//m.wg.Done()
//		return
//	}
//	user <- m.selectByComponent(findByEmail, command.Email.Address(), command.FindUserCommand.Password)
//	close(user)
//	return
//	//m.wg.Done()
//}
//
//func (m *MySqlUserRepository) FindUserByAlias(command domain.FindByAliasCommand, user chan *domain.User) {
//	const findByAlias = "SELECT * FROM users WHERE alias = ?;"
//
//	if err := m.newTrx(command.FindUserCommand.Log); err != nil {
//		user <- nil
//		close(user)
//		m.errorHandler()
//		//m.wg.Done()
//		return
//	}
//	_ = m.selectByComponent(findByAlias, command.Alias.Alias(), command.FindUserCommand.Password)
//	user <- &domain.User{}
//	//m.wg.Done()
//	close(user)
//	return
//}
//
//func (m *MySqlUserRepository) selectByComponent(query string, arg any, password string) *domain.User {
//	//defer func() {
//	//	m.wg.Done()
//	//}()
//	user := &UserSchema{}
//	_ = m.trx.QueryRow(query, arg).Scan(&user.Uuid, &user.Alias, &user.Name, &user.SecondName, &user.Email, &user.Password)
//
//	if user.Uuid == "" {
//		return nil
//	}
//
//	return domain.NewUser(domain.NewUserCommand{
//		Uuid:       user.Uuid,
//		Alias:      user.Alias,
//		Name:       user.Name,
//		SecondName: user.SecondName,
//		Email:      user.Email,
//		Password:   password,
//	}, m.log)
//}

//func (m *MySqlUserRepository) trxPrepare(log *returnLog.ReturnLog) error {
//	m.log = log
//	m.newTrx()
//	return m.err
//}

func (m *MySqlUserRepository) newTrx() error {
	trx, err := m.mySqlRepository.NewTrx()
	if err != nil {
		return err
	}

	m.trx = trx
	return nil
}

//func (m *MySqlUserRepository) errorHandler() {
//	if m.trx != nil {
//		_ = m.trx.Rollback()
//		//if rollbackErr := m.trx.Rollback(); rollbackErr != nil {
//		//	m.err = rollbackErr
//		//}
//	}
//
//	m.log.LogError(returnLog.NewErrorCommand{
//		Error:  m.err,
//		Caller: 3,
//	})
//	//m.wg.Done()
//	return
//}
