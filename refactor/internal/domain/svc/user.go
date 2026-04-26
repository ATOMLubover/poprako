package svc

import (
	"poprako-s/internal/domain/model/aggr"
	event_impl "poprako-s/internal/domain/model/event"
	"poprako-s/pkg/util"

	"golang.org/x/crypto/bcrypt"
)

type UserSvc struct{}

func NewUserSvc() UserSvc {
	return UserSvc{}
}

func (UserSvc) NewUserReg(inv *aggr.MemberInv, name string, pwd string) (*aggr.UserReg, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	uid := util.GenId("user")

	reg := &aggr.UserReg{
		Id:       uid,
		Qid:      inv.InviteeQid,
		Nickname: name,
		PwdHash:  string(hash),
	}

	reg.PushEv(event_impl.NewUserRegEv(inv.InvitorId, inv.InviteeQid, inv.TeamId))

	return reg, nil
}
