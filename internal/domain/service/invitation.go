package service

import "github.com/google/uuid"

func GenerateInvitationCode() (string, error) {
	// 为简便起见，裁剪并取生成的 UUID 的后 8 位作为邀请码
	uuid, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	return uuid.String()[len(uuid)-8:], nil
}
