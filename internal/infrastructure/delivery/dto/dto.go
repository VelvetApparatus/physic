package dto

import (
	"github.com/google/uuid"
	collectionDto "physk/internal/domain/aggregates/collection/dto"
	userDTO "physk/internal/domain/aggregates/user/dto"
)

type UserCreateRequest struct {
	Role     string `json:"role"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (u *UserCreateRequest) ToAggregateDTO() userDTO.CreateUserDTO {
	return userDTO.CreateUserDTO{
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

func (u *UserLoginRequest) ToAggregateDTO() userDTO.LoginUserDTO {
	return userDTO.LoginUserDTO{
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

type CreateCollectionRequest struct {
	Name string `json:"name"`
}

func (c *CreateCollectionRequest) ToAggregateDTO() collectionDto.CreateCollection {
	return collectionDto.CreateCollection{Name: c.Name}
}

type CreateCollectionResponse struct {
	ID uuid.UUID `json:"id"`
}

type AddImageToCollectionRequest struct {
	CollectionID uuid.UUID
	Name         string
	ContentType  string
	Data         []byte
}

func (a *AddImageToCollectionRequest) ToAggregateDTO() collectionDto.AddImageToCollection {
	return collectionDto.AddImageToCollection{
		CollectionID: a.CollectionID,
		Name:         a.Name,
		ContentType:  a.ContentType,
		Data:         a.Data,
	}
}

type AddImageToCollectionResponse struct {
	ID uuid.UUID `json:"id"`
}

type DeleteCollectionRequest struct {
	CollectionID uuid.UUID `json:"collection_id"`
}

func (d *DeleteCollectionRequest) ToAggregateDTO() collectionDto.DeleteCollection {
	return collectionDto.DeleteCollection{CollectionID: d.CollectionID}
}

type GetImageIdsByCollectionIDRequest struct {
	CollectionID uuid.UUID `json:"collection_id"`
}

type GetImagesIDsByCollectionIDResponse struct {
	ImageIDs []uuid.UUID `json:"image_ids"`
}

type GetImageByID struct {
	ID uuid.UUID `json:"id"`
}
