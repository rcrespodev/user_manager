package domain

import (
	"bytes"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"html/template"
	"path/filepath"
)

type WelcomeEmail struct {
	subject string
	user    *domain.User
	to      string
	body    []byte
	log     *returnLog.ReturnLog
}

func NewWelcomeEmail(user *domain.User, templatePath string, log *returnLog.ReturnLog) *WelcomeEmail {
	const (
		subject = "Welcome to User Manager App\n"
	)

	w := &WelcomeEmail{
		subject: subject,
		user:    user,
		to:      user.Email().Address(),
		log:     log,
		body:    nil,
	}

	w.newWelcomeBody(templatePath)
	if log.Error() != nil {
		return nil
	}

	return w
}

func (w *WelcomeEmail) newWelcomeBody(templatePath string) {
	templatePath, err := filepath.Abs(templatePath)
	if err != nil {
		w.log.LogError(returnLog.NewErrorCommand{Error: err})
		return
	}

	t, err := template.ParseFiles(templatePath)
	if err != nil {
		w.log.LogError(returnLog.NewErrorCommand{Error: err})
		return
	}

	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, w.bodyVars()); err != nil {
		w.log.LogError(returnLog.NewErrorCommand{Error: err})
		return
	}

	w.body = buffer.Bytes()
}

func (w *WelcomeEmail) bodyVars() interface{} {
	return struct {
		UserAlias string
	}{
		UserAlias: w.user.Alias().Alias(),
	}
}

func (w *WelcomeEmail) Subject() string {
	return w.subject
}

func (w *WelcomeEmail) To() string {
	return w.to
}

func (w *WelcomeEmail) Body() []byte {
	return w.body
}
