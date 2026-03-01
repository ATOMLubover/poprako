package service

import "github.com/google/uuid"

type InvitationService interface {
	GenerateInvitationCode() (string, error)
}

type invitationService struct{}

func NewInvitationService() InvitationService {
	return &invitationService{}
}

func (s *invitationService) GenerateInvitationCode() (string, error) {
	// 为简便起见，裁剪并取生成的 UUID 的后 8 位作为邀请码
	uuid, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	return uuid.String()[len(uuid)-8:], nil
}
