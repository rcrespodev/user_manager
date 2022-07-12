package mySql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rcrespodev/user_manager/pkg/kernel/repository/mySql/schemas"
	"log"
	"os"
)

type MySqlRepository struct {
	mySqlClient *sql.DB
}

const driver = "mysql"

func NewMySqlRepository(mySqlClient *sql.DB) *MySqlRepository {
	mySqlRepository := &MySqlRepository{}

	switch mySqlClient {
	case nil:
		user := os.Getenv("MYSQL_USER")
		if user == "" {
			log.Fatalf("env %v not found", "MYSQL_USER")
		}
		password := os.Getenv("MYSQL_PASSWORD")
		if password == "" {
			log.Fatalf("env %v not found", "MYSQL_PASSWORD")
		}

		userManagerDb := os.Getenv("MYSQL_DATABASE")
		if userManagerDb == "" {
			log.Fatalf("env %v not found", "MYSQL_DATABASE")
		}

		port := os.Getenv("MYSQL_PORT")
		if port == "" {
			log.Fatalf("env %v not found", "MYSQL_PORT")
		}

		//docker network
		db, err := sql.Open(driver, fmt.Sprintf("%v:%v@tcp(mysql:3306)/%v", user, password, userManagerDb))
		//dockertest
		//db, err := sql.Open("mysql", os.Getenv("DATASOURCE"))
		//db, err := sql.Open(driver, fmt.Sprintf("%v:%v@(localhost:%v)/%v", user, password, port, userManagerDb))
		if err != nil {
			log.Fatalf("Mysql %v", err)
		}

		if err = db.Ping(); err != nil {
			log.Fatal(err)
		}

	default:
		mySqlRepository.mySqlClient = mySqlClient
	}

	_, err := mySqlRepository.mySqlClient.Exec(schemas.UserSchema)
	if err != nil {
		log.Fatal(err)
	}

	return mySqlRepository
}

func (m MySqlRepository) NewTrx() (*sql.Tx, error) {
	return m.mySqlClient.BeginTx(context.Background(), nil)
}
