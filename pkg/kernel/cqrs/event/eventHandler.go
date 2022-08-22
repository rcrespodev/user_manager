package event

type Handler interface {
	Handle(event Event)
}
