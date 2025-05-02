package controller

import (
	"physk/internal/services/token_service"
	"physk/internal/usecase"
)

type Controller struct {
	write usecase.WriteModel
	read  usecase.ReadModel

	tokenService token_service.TokenService
}
