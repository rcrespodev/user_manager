package repository

import (
	"database/sql"
	emailSenderDomain "github.com/rcrespodev/user_manager/pkg/app/emailSender/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/repository/mySql"
	"time"
)

type MySqlSentEmailRepository struct {
	mySqlRepository *mySql.MySqlRepository
}

func NewMySqlSentEmailRepository(mySqlClient *sql.DB) *MySqlSentEmailRepository {
	return &MySqlSentEmailRepository{
		mySqlRepository: mySql.NewMySqlRepository(mySqlClient),
	}
}

func (m *MySqlSentEmailRepository) Save(schema *emailSenderDomain.SentEmailSchema, log *returnLog.ReturnLog) {
	const (
		saveSentEmail = `
		INSERT INTO sent_email (user_uuid, sent, sent_on, error) VALUES (?, ?, ?, ?);
		`
	)

	trx, err := m.mySqlRepository.NewTrx()
	if err != nil {
		log.LogError(returnLog.NewErrorCommand{Error: err})
		return
	}

	result, err := trx.Exec(saveSentEmail, schema.UserUuid, schema.Sent,
		schema.SentOn.Format("2006-01-02 15:04:05"), schema.Error)

	if err != nil {
		_ = trx.Rollback()
		log.LogError(returnLog.NewErrorCommand{Error: err})
		return
	}
	if rows, err := result.RowsAffected(); rows == 0 || err != nil {
		log.LogError(returnLog.NewErrorCommand{Error: err})
	}

	if err := trx.Commit(); err != nil {
		_ = trx.Rollback()
		log.LogError(returnLog.NewErrorCommand{Error: err})
		return
	}

	return
}

func (m *MySqlSentEmailRepository) Get(userUuid string) []*emailSenderDomain.SentEmailSchema {
	var schemas []*emailSenderDomain.SentEmailSchema
	const (
		getSentEmail = `
			SELECT s.user_uuid, s.sent, s.sent_on, s.error
			FROM sent_email as s
			WHERE user_uuid = ?;
		`
	)

	trx, err := m.mySqlRepository.NewTrx()
	if err != nil {
		return nil
	}
	rows, err := trx.Query(getSentEmail, userUuid)
	if err == sql.ErrNoRows {
		return nil
	}

	if rows == nil {
		return nil
	}

	if rows.Err() != nil {
		return nil
	}

	for rows.Next() {
		var (
			//schema *emailSenderDomain.SentEmailSchema
			uuid   string
			sent   bool
			sentOn string
			e      string
		)

		//var schema *emailSenderDomain.SentEmailSchema
		if err = rows.Scan(&uuid, &sent, &sentOn, &e); err != nil {
			continue
		}

		sentOnTime, err := time.Parse("2006-01-02 15:04:05", sentOn)
		if err != nil {
			continue
		}
		schema := &emailSenderDomain.SentEmailSchema{
			UserUuid: uuid,
			Sent:     sent,
			SentOn:   sentOnTime,
			Error:    e,
		}

		schemas = append(schemas, schema)
	}

	return schemas
}
