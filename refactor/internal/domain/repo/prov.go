package repo_iface

// `Prov` tries to retrieve repositories from a given context.
// NOTE: it panics if the context is missing any of the required repositories, so it should
// only be constructed by `TxnCtrl`.
type Prov interface {
	UserRepo() UserRepo
	TeamRepo() TeamRepo
	MemberRepo() MemberRepo
	MemberInvRepo() MemberInvRepo
	WorksetRepo() WorksetRepo
	ComicRepo() ComicRepo
	ChapterRepo() ChapterRepo
	PageRepo() PageRepo
	AssignmentInvRepo() AssignmentInvRepo
	AssignmentRepo() AssignmentRepo
	UserStatsRepo() UserStatsRepo
	OssMsgRepo() OssMsgRepo
}
