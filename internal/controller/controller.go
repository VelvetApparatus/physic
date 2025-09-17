package controller

import (
	tokenServ "physk/internal/services/token_service"
	"physk/internal/usecase"
)

type Controller struct {
	write usecase.WriteModel
	read  usecase.ReadModel

	tokenService tokenServ.TokenService
}

func NewController(
	writeModel usecase.WriteModel,
	readModel usecase.ReadModel,
) *Controller {
	return &Controller{
		write:        writeModel,
		read:         readModel,
		tokenService: tokenServ.NewTokenService(),
	}
}
