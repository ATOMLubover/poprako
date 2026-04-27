package event

import event_iface "poprako-s/internal/event"

// Event type identifiers for user-domain events
const (
	// `EvUserLogin` is emitted after a successful password login
	EvUserLogin event_iface.EvTyp = "event:user_login"
	// `EvUserReg` is emitted after a new user registration completes
	EvUserReg event_iface.EvTyp = "event:user_reg"
)
