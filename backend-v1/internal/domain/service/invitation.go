package service

import "github.com/google/uuid"

func GenerateInvitationCode() (string, error) {
	// 为简便起见，裁剪并取生成的 UUID 的后 8 位作为邀请码
	rawUUID, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	uuidString := rawUUID.String()

	return uuidString[len(uuidString)-6:], nil
}
