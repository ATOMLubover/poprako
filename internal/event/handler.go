package event_iface

import "context"

// `EvHandler` handles a dispatched domain event.
type EvHandler interface {
	// `EvTyp` returns an identifier of target event.
	EvTyp() EvTyp

	// `Handle` processes a domain event in a specific context.
	// Though handlers work often in async context, some infomation
	// that is not contained in `ev` can be passed within it.
	//
	// NOTE: As no error from handler can be handled, `Handle`
	// returns void.
	Handle(cx context.Context, ev Event)
}
