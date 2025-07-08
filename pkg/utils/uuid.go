package utils

import "github.com/google/uuid"

func (u *utility) UUIDGenerate() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
