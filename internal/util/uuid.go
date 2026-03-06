package util

import "github.com/google/uuid"

func GenerateUUID() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.NewString()
	}

	return id.String()
}
