package event

import event_iface "poprako-s/internal/event"

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
