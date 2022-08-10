package mySql

import (
	"context"
	"database/sql"
	_ "database/sql/driver"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/config"
	"github.com/rcrespodev/user_manager/pkg/kernel/repository/mySql/schemas"
	"log"
	_ "os"
)

type MySqlRepository struct {
	mySqlClient *sql.DB
}

const driver = "mysql"

func NewMySqlRepository(mySqlClient *sql.DB) *MySqlRepository {
	mySqlRepository := &MySqlRepository{}

	switch mySqlClient {
	case nil:
		mySqlRepository.mySqlClient = mySqlRepository.newConnection()
	default:
		mySqlRepository.mySqlClient = mySqlClient
	}

	_, err := mySqlRepository.mySqlClient.Exec(schemas.UserSchema)
	if err != nil {
		log.Fatal(err)
	}

	return mySqlRepository
}

func (m *MySqlRepository) NewTrx() (*sql.Tx, error) {
	return m.mySqlClient.BeginTx(context.Background(), nil)
}

func (m *MySqlRepository) newConnection() *sql.DB {
	//docker network
	mySqlConf := config.Conf.MySql
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=5s",
		mySqlConf.User, mySqlConf.Password, mySqlConf.Host, mySqlConf.Port, mySqlConf.Database)

	db, err := sql.Open(driver, dataSource)
	if err != nil {
		log.Fatalf("Mysql connnection %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}
