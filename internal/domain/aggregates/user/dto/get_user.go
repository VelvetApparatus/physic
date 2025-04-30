package dto

import (
	"github.com/google/uuid"
	vo "physk/internal/domain/aggregates/user/value_objects"
)

type GetUserDTO struct {
	ID           uuid.UUID
	Role         vo.UserRole
	Login        string
	PasswordHash string
	Username     string
	Email        string
}
