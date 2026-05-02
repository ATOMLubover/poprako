package svc

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/event"
	"poprako-s/pkg/util"

	"golang.org/x/crypto/bcrypt"
)

// `UserSvc` provides domain services for user creation and credential management
type UserSvc struct{}

// `NewUserSvc` creates a new `UserSvc`
func NewUserSvc() UserSvc {
	return UserSvc{}
}

// `NewUserReg` builds a `UserReg` aggregate from an invitation and registration inputs,
// hashing the password and pushing a `UserRegEv` onto the aggregate
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

	reg.PushEv(event.NewUserRegEv(inv.InvitorId, inv.InviteeQid, inv.TeamId))

	return reg, nil
}
