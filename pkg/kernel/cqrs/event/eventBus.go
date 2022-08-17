package event

type Bus interface {
	Subscribe(events []Event)
	Publish(eventId Id, handler Handler)
}
