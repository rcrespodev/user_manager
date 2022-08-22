package event

type MockEventBus struct {
}

func (m MockEventBus) Subscribe(eventId Id, handler Handler) {
	return
}

func (m MockEventBus) Publish(events []Event) {
	return
}
