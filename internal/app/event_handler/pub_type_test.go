package event_handler

import (
	"testing"

	"poprako-s/internal/domain/event"
)

func TestTxnBoundHandlersUseSyncPubType(t *testing.T) {
	handlers := []event.EventHandler{
		NewUnitSaveHandler(),
		NewAssignmentCreateHandler(),
		NewAssignmentRemoveHandler(),
		NewComicCreateHandler(),
		NewComicRemoveHandler(),
		NewChapterCreateHandler(),
		NewChapterRemoveHandler(),
		NewChapterPublishedHandler(),
		NewChapterCreatorAssignedHandler(),
	}

	for _, handler := range handlers {
		if handler.PubType() != event.PubTypeSync {
			t.Fatalf("handler %T must use PubTypeSync", handler)
		}
	}
}
