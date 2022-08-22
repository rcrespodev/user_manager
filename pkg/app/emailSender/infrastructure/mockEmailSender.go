package infrastructure

import (
	emailSenderDomain "github.com/rcrespodev/user_manager/pkg/app/emailSender/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type MockEmailSender struct {
}

func (m MockEmailSender) SendEmail(command *emailSenderDomain.SendEmailCommand, log *returnLog.ReturnLog) {
	return
}
