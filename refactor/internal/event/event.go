package event_iface

// `Event` represents a domain event.
type Event interface {
	// `EvTyp` returns an identifier of this event.
	// A `EventHandler` finds its corresponding evnets
	// according to this identifier.
	EvTyp() EvTyp

	// `Payload` returns the original type of the implementation
	// of the `Event` interface.
	Payload() any
}
