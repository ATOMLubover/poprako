package app_iface

import (
	"context"

	"poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
)

// `AssignmentApp` defines use-cases for assignment.
type AssignmentApp interface {
	// `ListByChapter` lists assignments under one chapter.
	ListByChapter(cx context.Context, currUid string, args *val.ListAssignmentByChapterArgs) app_res.AppRes[[]val.AssignmentVal]

	// `ListByUser` lists all assignments of current user.
	ListByUser(cx context.Context, currUid string, args *val.ListMyAssignmentArgs) app_res.AppRes[[]val.AssignmentVal]

	// `Upsert` executes put-semantics upsert for assignment roles.
	// If role mask is zero, this use-case redirects to delete semantics.
	Upsert(cx context.Context, currUid string, args *val.UpsertAssignmentArgs) app_res.AppRes[app_res.None]

	// `Delete` deletes one assignment by id.
	Delete(cx context.Context, currUid string, assignmentId string) app_res.AppRes[app_res.None]
}
