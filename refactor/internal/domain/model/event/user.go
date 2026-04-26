package event_impl

import event_iface "poprako-s/internal/event"

type UserLoginEv struct {
	UserId string
}

func NewUserLoginEv(userId string) event_iface.Event {
	return &UserLoginEv{UserId: userId}
}

func (e *UserLoginEv) EvTyp() event_iface.EvTyp {
	return EvUserLogin
}

func (e *UserLoginEv) Payload() any {
	return e
}

type UserRegEv struct {
	InvitorId  string
	InviteeQid string
	TeamId     string
}

func NewUserRegEv(invitorId string, inviteeQid string, teamId string) event_iface.Event {
	return &UserRegEv{
		InvitorId:  invitorId,
		InviteeQid: inviteeQid,
		TeamId:     teamId,
	}
}

func (e *UserRegEv) EvTyp() event_iface.EvTyp {
	return EvUserReg
}

func (e *UserRegEv) Payload() any {
	return e
}
