package infrastructure

import (
	"fmt"
	emailSenderDomain "github.com/rcrespodev/user_manager/pkg/app/emailSender/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"net/smtp"
)

type SmtpEmailSender struct {
	userName  string
	host      string
	port      string
	emailAuth smtp.Auth
}

type EmailAuthConf struct {
	Host     string
	From     string
	Password string
	Port     string
}

func NewSmtpEmailSender(conf EmailAuthConf) *SmtpEmailSender {
	return &SmtpEmailSender{
		emailAuth: smtp.PlainAuth("", conf.From, conf.Password, conf.Host),
		userName:  conf.From,
	}
}

func (s SmtpEmailSender) SendEmail(command *emailSenderDomain.SendEmailCommand, log *returnLog.ReturnLog) {
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	message := []byte(fmt.Sprintf("Subject: %s\n%s\n%s\n", command.Subject, mime, string(command.Body)))
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	to := make([]string, 1)
	to[0] = command.To
	if err := smtp.SendMail(addr, s.emailAuth, s.userName, to, message); err != nil {
		log.LogError(returnLog.NewErrorCommand{Error: err})
	}
}
