package mock_event

import (
	"context"

	"poprako-s/internal/domain/event"
)

type EventBus struct {
	handlers      map[event.EventType][]event.EventHandler
	pubCalls      [][]event.Event
	pubAsyncCalls [][]event.Event
	pubErr        error
	pubAsyncErr   error
}

var _ event.EventBus = (*EventBus)(nil)

func NewMockEventBus() event.EventBus {
	return &EventBus{
		handlers: make(map[event.EventType][]event.EventHandler),
	}
}

func (b *EventBus) SetPubErr(err error) {
	b.pubErr = err
}

func (b *EventBus) SetPubAsyncErr(err error) {
	b.pubAsyncErr = err
}

func (b *EventBus) PubCalls() [][]event.Event {
	result := make([][]event.Event, 0, len(b.pubCalls))

	for _, call := range b.pubCalls {
		result = append(result, append([]event.Event(nil), call...))
	}

	return result
}

func (b *EventBus) PubAsyncCalls() [][]event.Event {
	result := make([][]event.Event, 0, len(b.pubAsyncCalls))

	for _, call := range b.pubAsyncCalls {
		result = append(result, append([]event.Event(nil), call...))
	}

	return result
}

func (b *EventBus) Pub(cx context.Context, ev []event.Event) error {
	b.pubCalls = append(b.pubCalls, append([]event.Event(nil), ev...))

	if b.pubErr != nil {
		return b.pubErr
	}

	// Execute PubTypeSync handlers inline to match real EventBus behavior
	for _, e := range ev {
		for _, h := range b.handlers[e.EventType()] {
			if h.PubType() == event.PubTypeSync {
				if err := h.Handle(cx, e); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (b *EventBus) PubAsync(cx context.Context, ev []event.Event) error {
	b.pubAsyncCalls = append(b.pubAsyncCalls, append([]event.Event(nil), ev...))

	return b.pubAsyncErr
}

func (b *EventBus) Sub(h event.EventHandler) error {
	ty := h.EventType()
	newHandlers := make([]event.EventHandler, len(b.handlers[ty])+1)
	copy(newHandlers, b.handlers[ty])
	newHandlers[len(newHandlers)-1] = h
	b.handlers[ty] = newHandlers
	return nil
}

func (b *EventBus) SubUnsafe(h event.EventHandler) error {
	return b.Sub(h)
}

func (b *EventBus) Unsub(h event.EventHandler) error {
	ty := h.EventType()
	handlers := b.handlers[ty]
	for i, handler := range handlers {
		if handler == h {
			b.handlers[ty] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
	return nil
}

func (b *EventBus) UnsubUnsafe(h event.EventHandler) error {
	return b.Unsub(h)
}

func (b *EventBus) Close() {}
