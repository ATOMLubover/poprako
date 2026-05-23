package repo_infra

import (
	repo_iface "poprako-s/internal/domain/repo"

	"gorm.io/gorm"
)

type provImpl struct {
	gdb *gorm.DB
}

func newProv(gdb *gorm.DB) repo_iface.Prov {
	return &provImpl{gdb: gdb}
}

func (p *provImpl) UserRepo() repo_iface.UserRepo {
	return NewUserRepo(p.gdb)
}

func (p *provImpl) TeamRepo() repo_iface.TeamRepo {
	return NewTeamRepo(p.gdb)
}

func (p *provImpl) MemberRepo() repo_iface.MemberRepo {
	return NewMemberRepo(p.gdb)
}

func (p *provImpl) AnnouncementRepo() repo_iface.AnnouncementRepo {
	return NewAnnouncementRepo(p.gdb)
}

func (p *provImpl) CommentRepo() repo_iface.CommentRepo {
	return NewCommentRepo(p.gdb)
}

func (p *provImpl) MemberInvRepo() repo_iface.MemberInvRepo {
	return NewMemberInvRepo(p.gdb)
}

func (p *provImpl) WorksetRepo() repo_iface.WorksetRepo {
	return NewWorksetRepo(p.gdb)
}

func (p *provImpl) ComicRepo() repo_iface.ComicRepo {
	return NewComicRepo(p.gdb)
}

func (p *provImpl) ChapterRepo() repo_iface.ChapterRepo {
	return NewChapterRepo(p.gdb)
}

func (p *provImpl) PageRepo() repo_iface.PageRepo {
	return NewPageRepo(p.gdb)
}

func (p *provImpl) UnitRepo() repo_iface.UnitRepo {
	return NewUnitRepo(p.gdb)
}

func (p *provImpl) AssignmentInvRepo() repo_iface.AssignmentInvRepo {
	return NewAssignmentInvRepo(p.gdb)
}

func (p *provImpl) AssignmentRepo() repo_iface.AssignmentRepo {
	return NewAssignmentRepo(p.gdb)
}

func (p *provImpl) OssMsgRepo() repo_iface.OssMsgRepo {
	return NewOssMsgRepo(p.gdb)
}
