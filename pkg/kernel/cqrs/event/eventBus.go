package event

type Bus interface {
	Subscribe(eventId Id, handler Handler)
	Publish(events []Event)
}
