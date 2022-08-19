package registerUser

import (
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/event"
)

type UserRegistered struct {
	baseEvent *event.BaseEvent
}

func NewUserRegistered(command Command) *UserRegistered {
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
