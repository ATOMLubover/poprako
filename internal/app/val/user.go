package val

import token_iface "poprako-s/internal/domain/ext/token"

type UserVal struct {
	Id string `json:"id"`

	Qid      string `json:"qid"`
	Nickname string `json:"nickname"`

	AvatarUrl      string `json:"avatar_url"`
	AvatarUploaded bool   `json:"avatar_uploaded"`

	IsSuperAdmin bool `json:"is_super_admin"`

	LastActiveAt int64 `json:"last_active_at"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

type UserLoginArgs struct {
	Qid string `json:"qid"`
	Pwd string `json:"password"`
}

type UserLoginRes struct {
	UserId string                  `json:"user_id"`
	Token  token_iface.SignedToken `json:"token"`
}

type UserRegArgs struct {
	Qid  string `json:"qid"`
	Name string `json:"name"`
	Pwd  string `json:"password"`

	InvCode string `json:"invitation_code"`
}

type UserRegRes struct {
	UserId string                  `json:"user_id"`
	Token  token_iface.SignedToken `json:"token"`
}

// `UserUpdArgs` holds mutable fields for one user put update.
type UserUpdArgs struct {
	Id string `json:"id"`

	Name string `json:"name"`
	Qid  string `json:"qid"`
}

// `ResvUserAvatarArgs` provides the file extension
// of the avatar to be uploaded, so that the server and
// generate a proper upload URL and save path.
type ResvUserAvatarArgs struct {
	UserId  string `json:"user_id"`
	FileExt string `json:"file_extension"`
}

type ResvUserAvatarRes struct {
	PutUrl string `json:"put_url"`
}

// `UserRel` holds two user ids, wich maybe useful when
// we need to check the relationship between two users, e.g. whether
// curr user is permitted to update the other one's info.
type UserRel struct {
	SubUid string `json:"-"`
	ObjUid string `json:"-"`
}
