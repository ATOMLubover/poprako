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

func (b *EvBase) clearEv() {
	b.ev = nil
}

func (b *EvBase) PullEv() []Event {
	ev := b.ev
	b.clearEv()

	return ev
}
