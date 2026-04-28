package app_iface

import (
	"context"

	"poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
)

// `AssignmentInvApp` defines use-cases for assignment invitation.
type AssignmentInvApp interface {
	// `ListByChapter` lists invitations under one chapter.
	ListByChapter(cx context.Context, currUid string, args *val.ListAssignmentInvArgs) res.AppRes[[]val.AssignmentInvVal]

	// `Create` creates one invitation for assignment.
	Create(cx context.Context, currUid string, args *val.CreateAssignmentInvArgs) res.AppRes[val.CreateAssignmentInvRes]

	// `Remove` removes one invitation by id.
	Remove(cx context.Context, currUid string, invId string) res.AppRes[res.None]

	// `JoinByInvCode` joins one chapter by invitation code.
	JoinByInvCode(cx context.Context, currUid string, args *val.JoinAssignmentInvArgs) res.AppRes[res.None]
}
