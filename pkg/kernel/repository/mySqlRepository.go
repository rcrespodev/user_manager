package repository

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

type MySqlRepository struct {
	mySqlClient *sql.DB
}

func NewMySqlRepository() *MySqlRepository {
	const driver = "mysql"

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

	//db, err := sql.Open(driver, fmt.Sprintf("%v:%v@tcp(0.0.0.0:3306)/%v", user, password, userManagerDb))
	db, err := sql.Open(driver, fmt.Sprintf("%v:%v@tcp(mysql:3306)/%v", user, password, userManagerDb))
	if err != nil {
		log.Fatalf("Mysql %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return &MySqlRepository{
		mySqlClient: db,
	}
}

func (m MySqlRepository) NewTrx() (*sql.Tx, error) {
	return m.mySqlClient.BeginTx(context.Background(), nil)
}
