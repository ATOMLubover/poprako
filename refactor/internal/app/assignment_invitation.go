package app_iface

import (
	"context"

	"poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
)

// `AssignmentInvApp` defines use-cases for assignment invitation.
type AssignmentInvApp interface {
	// `ListByChapter` lists invitations under one chapter.
	ListByChapter(cx context.Context, currUid string, args *val.ListAssignmentInvArgs) app_res.AppRes[[]val.AssignmentInvVal]

	// `Create` creates one invitation for assignment.
	Create(cx context.Context, currUid string, args *val.CreateAssignmentInvArgs) app_res.AppRes[val.CreateAssignmentInvRes]

	// `Delete` deletes one invitation by id.
	Delete(cx context.Context, currUid string, invId string) app_res.AppRes[app_res.None]

	// `JoinByInvCode` joins one chapter by invitation code.
	JoinByInvCode(cx context.Context, currUid string, args *val.JoinAssignmentInvArgs) app_res.AppRes[app_res.None]
}
