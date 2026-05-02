package event_iface

// `EvBase` wraps a list of events.
// It provides some simple methods to push and pull
// events stored in it.
type EvBase struct {
	ev []Event
}

// `PushEv` append one event in the pending list of `EventBase`.
func (b *EvBase) PushEv(e Event) {
	b.ev = append(b.ev, e)
}

// `clearEv` clears all events in pending list of `EventBase`.
// In case of memory leak, never ref **SINGLE** element of
// b.ev with raw calls.
func (b *EvBase) clearEv() {
	clear(b.ev)
}

func (b *EvBase) PullEv() []Event {
	ev := b.ev
	b.clearEv()

	return ev
}
