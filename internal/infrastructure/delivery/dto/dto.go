package dto

import (
	"github.com/google/uuid"
	"physk/internal/domain/aggregates/user/dto"
)

type UserCreateRequest struct {
	Role     string
	Login    string
	Password string
	Username string
	Email    string
}

func (u UserCreateRequest) ToAggregateDTO() dto.CreateUserDTO {
	return dto.CreateUserDTO{
		Role:     u.Role,
		Login:    u.Login,
		Password: u.Password,
		Username: u.Username,
		Email:    u.Email,
	}
}

type UserCreateResponse struct {
	ID       uuid.UUID `json:"id"`
	Role     string    `json:"role"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (u UserLoginRequest) ToAggregateDTO() dto.LoginUserDTO {
	return dto.LoginUserDTO{
		Login:    u.Login,
		Password: u.Password,
	}
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type GetMeResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
}
