package domain

import returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"

type EmailSender interface {
	SendEmail(command *SendEmailCommand, log *returnLog.ReturnLog)
}

type SendEmailCommand struct {
	To      string
	Subject string
	Body    []byte
}
