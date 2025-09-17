package dto

import (
	"github.com/google/uuid"
)

type GetUserDTO struct {
	ID           uuid.UUID
	Role         string
	Login        string
	PasswordHash string
	Username     string
	Email        string
}
