package value_objects

import (
	"physk/pkg/hasher"
)

const (
	passwordMinLen = 8
	passwordMaxLen = 40
)

type UserPasswordHash string

func (uph UserPasswordHash) String() string { return string(uph) }

// todo: what about errors ?

func ParsePassword(password string) (res UserPasswordHash, valid bool) {
	if len(password) < passwordMinLen || len(password) > passwordMaxLen {
		return "", false
	}

	passwordHash, err := hasher.HashString(password)
	if err != nil {
		return "", false
	}

	return UserPasswordHash(passwordHash), true
}
