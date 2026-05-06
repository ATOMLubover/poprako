package event_infra

import (
	"context"

	event_iface "poprako-s/internal/event"

	"go.uber.org/zap"
)

type evBusImpl struct {
	evHdl map[event_iface.EvTyp][]event_iface.EvHandler
	tx    chan *workUnit
	done  chan struct{}
}

func NewEvBus() event_iface.EvBus {
	const TX_SIZE = 1024

	evHdl := make(map[event_iface.EvTyp][]event_iface.EvHandler)
	tx := make(chan *workUnit, TX_SIZE) // Buffer size can be adjusted based on expected load.
	done := make(chan struct{})

	return &evBusImpl{
		evHdl: evHdl,
		tx:    tx,
		done:  done,
	}
}

func (b *evBusImpl) Pub(cx context.Context, ev []event_iface.Event) {
	wu := &workUnit{
		cx: cx,
		ev: ev,
	}

	// Non-blocking transmission.
	select {
	case b.tx <- wu:
		zap.L().Debug(
			"[evBusImpl.Pub] published event to bus",
			zap.Any("event", ev),
		)
	default:
		zap.L().Error(
			"[evBusImpl.Pub] failed to publish event to bus: bus is busy",
			zap.Any("event", ev),
		)
	}
}

func (b *evBusImpl) Sub(h event_iface.EvHandler) {
	typ := h.EvTyp()

	if _, exists := b.evHdl[typ]; !exists {
		b.evHdl[typ] = []event_iface.EvHandler{}
	}

	b.evHdl[typ] = append(b.evHdl[typ], h)

	zap.L().Debug(
		"[evBusImpl.Sub] subscribed handler to event type",
		zap.Any("event_type", typ),
	)
}

func (b *evBusImpl) Run() {
	rx := b.tx
	done := b.done

	onWu := func(wu *workUnit) {
		for _, ev := range wu.ev {
			handlers, exists := b.evHdl[ev.EvTyp()]
			if !exists {
				zap.L().Warn(
					"[evBusImpl.Run] no handlers for event type",
					zap.Any("event_type", ev.EvTyp()),
				)
				continue
			}

			for _, h := range handlers {
				go h.Handle(wu.cx, ev)
			}
		}
	}

	go func() {
		for {
			select {
			case wu := <-rx:
				onWu(wu)
			case <-done:
				zap.L().Info("[evBusImpl.Run] event bus is closing")
				return
			}
		}
	}()
}

func (b *evBusImpl) Close() {
	b.done <- struct{}{}
}
