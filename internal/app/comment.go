package app_iface

import (
	"context"

	app_res "poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
)

// `CommentApp` defines application use-cases for team comments.
type CommentApp interface {
	// `List` lists team board comments.
	List(cx context.Context, currUid string, args *val.ListCommentArgs) app_res.AppRes[[]val.CommentVal]

	// `Create` creates one team board comment.
	Create(cx context.Context, currUid string, args *val.CreateCommentArgs) app_res.AppRes[val.CommentCreatedRes]
}
