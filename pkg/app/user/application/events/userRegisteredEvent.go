package events

import (
	"github.com/rcrespodev/user_manager/pkg/app/user/application/commands/register"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/event"
)

type UserRegistered struct {
	baseEvent *event.BaseEvent
}

func NewUserRegistered(command register.RegisterUserCommand) *UserRegistered {
	return &UserRegistered{
		baseEvent: event.NewBaseEvent(event.NewBaseEventCommand{
			AggregateId: command.BaseCommand().AggregateId(),
			EventId:     event.UserRegistered,
			CommandUuid: command.BaseCommand().CommandUuid(),
		}),
	}
}

func (u UserRegistered) BaseEvent() *event.BaseEvent {
	return u.baseEvent
}
