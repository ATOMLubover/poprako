package mock_event

import "poprako-s/internal/domain/event"

type EventBus struct {
	pubCalls      [][]event.Event
	pubAsyncCalls [][]event.Event
	pubErr        error
	pubAsyncErr   error
}

var _ event.EventBus = (*EventBus)(nil)

func NewMockEventBus() event.EventBus {
	return &EventBus{}
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

func (b *EventBus) Pub(ev []event.Event) error {
	b.pubCalls = append(b.pubCalls, append([]event.Event(nil), ev...))
	return b.pubErr
}

func (b *EventBus) PubAsync(ev []event.Event) error {
	b.pubAsyncCalls = append(b.pubAsyncCalls, append([]event.Event(nil), ev...))
	return b.pubAsyncErr
}
func (b *EventBus) Sub(h event.EventHandler) error         { return nil }
func (b *EventBus) SubUnsafe(h event.EventHandler) error   { return nil }
func (b *EventBus) Unsub(h event.EventHandler) error       { return nil }
func (b *EventBus) UnsubUnsafe(h event.EventHandler) error { return nil }
func (b *EventBus) Close()                                 {}
