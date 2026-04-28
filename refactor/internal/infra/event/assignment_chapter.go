package event_infra

import (
	"context"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/event"
	repo_iface "poprako-s/internal/domain/repo"
	event_iface "poprako-s/internal/event"

	"go.uber.org/zap"
)

// `UpdateUserStatsOnAssignmentCreatedHandler` handles assignment-created side effects.
type UpdateUserStatsOnAssignmentCreatedHandler struct {
	userStatsRepo repo_iface.UserStatsRepo
}

// `NewUpdateUserStatsOnAssignmentCreatedHandler` creates one assignment-created handler.
func NewUpdateUserStatsOnAssignmentCreatedHandler(userStatsRepo repo_iface.UserStatsRepo) *UpdateUserStatsOnAssignmentCreatedHandler {
	return &UpdateUserStatsOnAssignmentCreatedHandler{userStatsRepo: userStatsRepo}
}

// `EvTyp` returns target event type.
func (h *UpdateUserStatsOnAssignmentCreatedHandler) EvTyp() event_iface.EvTyp {
	return event.EvAssignmentCreated
}

// `Handle` applies user stats delta for assignment creation.
func (h *UpdateUserStatsOnAssignmentCreatedHandler) Handle(_ context.Context, ev event_iface.Event) {
	payload, ok := ev.Payload().(*event.AssignmentCreatedEv)
	if !ok || payload == nil {
		zap.L().Error(
			"[UpdateUserStatsOnAssignmentCreatedHandler.Handle] invalid event payload",
			zap.Any("payload", ev.Payload()),
		)

		return
	}

	if err := h.userStatsRepo.Patch(&aggr.UserStatsPatch{
		UserId:                  payload.UserId,
		TotalAssignmentDelta:    1,
		ActiveAssignmentDelta:   1,
		FinishedAssignmentDelta: 0,
	}); err != nil {
		zap.L().Error(
			"[UpdateUserStatsOnAssignmentCreatedHandler.Handle] failed to patch user stats",
			zap.String("user_id", payload.UserId),
			zap.Error(err),
		)

		return
	}
}

// `UpdateUserStatsOnAssignmentRemovedHandler` handles assignment-removed side effects.
type UpdateUserStatsOnAssignmentRemovedHandler struct {
	userStatsRepo repo_iface.UserStatsRepo
}

// `NewUpdateUserStatsOnAssignmentRemovedHandler` creates one assignment-removed handler.
func NewUpdateUserStatsOnAssignmentRemovedHandler(userStatsRepo repo_iface.UserStatsRepo) *UpdateUserStatsOnAssignmentRemovedHandler {
	return &UpdateUserStatsOnAssignmentRemovedHandler{userStatsRepo: userStatsRepo}
}

// `EvTyp` returns target event type.
func (h *UpdateUserStatsOnAssignmentRemovedHandler) EvTyp() event_iface.EvTyp {
	return event.EvAssignmentRemoved
}

// `Handle` applies user stats delta for assignment removal.
func (h *UpdateUserStatsOnAssignmentRemovedHandler) Handle(_ context.Context, ev event_iface.Event) {
	payload, ok := ev.Payload().(*event.AssignmentRemovedEv)
	if !ok || payload == nil {
		zap.L().Error(
			"[UpdateUserStatsOnAssignmentRemovedHandler.Handle] invalid event payload",
			zap.Any("payload", ev.Payload()),
		)

		return
	}

	if payload.WasPublished {
		return
	}

	if err := h.userStatsRepo.Patch(&aggr.UserStatsPatch{
		UserId:                payload.UserId,
		ActiveAssignmentDelta: -1,
	}); err != nil {
		zap.L().Error(
			"[UpdateUserStatsOnAssignmentRemovedHandler.Handle] failed to patch user stats",
			zap.String("user_id", payload.UserId),
			zap.Error(err),
		)

		return
	}
}

// `UpdateUserStatsOnChapterPublishedHandler` handles chapter-published side effects.
type UpdateUserStatsOnChapterPublishedHandler struct {
	userStatsRepo repo_iface.UserStatsRepo
}

// `NewUpdateUserStatsOnChapterPublishedHandler` creates one chapter-published handler.
func NewUpdateUserStatsOnChapterPublishedHandler(userStatsRepo repo_iface.UserStatsRepo) *UpdateUserStatsOnChapterPublishedHandler {
	return &UpdateUserStatsOnChapterPublishedHandler{userStatsRepo: userStatsRepo}
}

// `EvTyp` returns target event type.
func (h *UpdateUserStatsOnChapterPublishedHandler) EvTyp() event_iface.EvTyp {
	return event.EvChapterPublished
}

// `Handle` converts active assignments to finished assignments on chapter publish.
func (h *UpdateUserStatsOnChapterPublishedHandler) Handle(_ context.Context, ev event_iface.Event) {
	payload, ok := ev.Payload().(*event.ChapterPublishedEv)
	if !ok || payload == nil {
		zap.L().Error(
			"[UpdateUserStatsOnChapterPublishedHandler.Handle] invalid event payload",
			zap.Any("payload", ev.Payload()),
		)

		return
	}

	for _, userId := range payload.AssignedUserIds {
		if userId == "" {
			continue
		}

		if err := h.userStatsRepo.Patch(&aggr.UserStatsPatch{
			UserId:                  userId,
			ActiveAssignmentDelta:   -1,
			FinishedAssignmentDelta: 1,
		}); err != nil {
			zap.L().Error(
				"[UpdateUserStatsOnChapterPublishedHandler.Handle] failed to patch user stats",
				zap.String("chapter_id", payload.ChapterId),
				zap.String("user_id", userId),
				zap.Error(err),
			)

			continue
		}
	}
}

// `UpdateUserStatsOnChapterRemovedHandler` handles chapter-removed side effects.
type UpdateUserStatsOnChapterRemovedHandler struct {
	userStatsRepo repo_iface.UserStatsRepo
}

// `NewUpdateUserStatsOnChapterRemovedHandler` creates one chapter-removed handler.
func NewUpdateUserStatsOnChapterRemovedHandler(userStatsRepo repo_iface.UserStatsRepo) *UpdateUserStatsOnChapterRemovedHandler {
	return &UpdateUserStatsOnChapterRemovedHandler{userStatsRepo: userStatsRepo}
}

// `EvTyp` returns target event type.
func (h *UpdateUserStatsOnChapterRemovedHandler) EvTyp() event_iface.EvTyp {
	return event.EvChapterRemoved
}

// `Handle` removes active assignment stats for chapter removal before publish.
func (h *UpdateUserStatsOnChapterRemovedHandler) Handle(_ context.Context, ev event_iface.Event) {
	payload, ok := ev.Payload().(*event.ChapterRemovedEv)
	if !ok || payload == nil {
		zap.L().Error(
			"[UpdateUserStatsOnChapterRemovedHandler.Handle] invalid event payload",
			zap.Any("payload", ev.Payload()),
		)

		return
	}

	if payload.WasPublished {
		return
	}

	for _, userId := range payload.AssignedUserIds {
		if userId == "" {
			continue
		}

		if err := h.userStatsRepo.Patch(&aggr.UserStatsPatch{
			UserId:                userId,
			ActiveAssignmentDelta: -1,
		}); err != nil {
			zap.L().Error(
				"[UpdateUserStatsOnChapterRemovedHandler.Handle] failed to patch user stats",
				zap.String("chapter_id", payload.ChapterId),
				zap.String("user_id", userId),
				zap.Error(err),
			)

			continue
		}
	}
}
