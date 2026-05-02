package event_iface

import "context"

// `EvBus` manages all events and corresponding handles.
// It also works as dispatcher.
type EvBus interface {
	// `Pub` publish a event on event bus.
	// All handlers subscribed to it will be notified
	// **asynchronously**.
	Pub(cx context.Context, ev []Event)

	// `Sub` subscribes a event handler to a
	// specific event type in **bootstrap stage**.
	// It may not be called during running stage.
	Sub(h EvHandler)

	// `Run` starts background dispatcher of event bus.
	Run()

	// `Close` closes background goroutine.
	Close()
}
