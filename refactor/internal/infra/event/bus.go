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
	evHdl := make(map[event_iface.EvTyp][]event_iface.EvHandler)
	tx := make(chan *workUnit, 100) // Buffer size can be adjusted based on expected load.
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
		zap.Any("handler", h),
	)
}

func (b *evBusImpl) Run() {
	rx := b.tx
	done := b.done

	go func() {
		for {
			select {
			case wu := <-rx:
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
			case <-done:
				zap.L().Info("[evBusImpl.Run] event bus is closing")
				break
			}
		}
	}()
}

func (b *evBusImpl) Close() {
	b.done <- struct{}{}
}
