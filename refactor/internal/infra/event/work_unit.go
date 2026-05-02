package event_infra

import (
	"context"

	event_iface "poprako-s/internal/event"
)

type workUnit struct {
	cx context.Context
	ev []event_iface.Event
}

func (u *workUnit) run() {
}
