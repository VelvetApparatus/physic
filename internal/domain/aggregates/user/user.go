package user

import (
	"fmt"
	"github.com/google/uuid"
	"physk/internal/domain/aggregates/user/context"
	userDto "physk/internal/domain/aggregates/user/dto"
	userErrors "physk/internal/domain/aggregates/user/errors"
	"physk/internal/domain/aggregates/user/rules"
	vo "physk/internal/domain/aggregates/user/value_objects"
	"physk/pkg/hasher"
)

type User struct {
	ID           uuid.UUID
	Role         vo.UserRole
	Login        string
	PasswordHash string
	Username     string
	Email        string
}

func GetUser(
	ctx context.UserCtx,
	dto userDto.GetUserDTO,
) (User, error) {

	// todo: validate too ?
	// i think -- yes

	aggregate := User{
		ID:           dto.ID,
		Role:         dto.Role,
		Login:        dto.Login,
		PasswordHash: dto.PasswordHash,
		Username:     dto.Username,
		Email:        dto.Email,
	}

	return aggregate, nil
}

func CreateUser(
	ctx context.UserCtx,
	dto userDto.CreateUserDTO,
) (User, error) {
	err := rules.CheckLoginUniqueness(ctx, dto.Login)
	if err != nil {
		return User{}, fmt.Errorf("check login uniqueness: %w", err)
	}

	err = rules.CheckEmailUniqueness(ctx, dto.Email)
	if err != nil {
		return User{}, fmt.Errorf("check email uniquenes: %w", err)
	}

	err = rules.CheckUsernameUniqueness(ctx, dto.Username)
	if err != nil {
		return User{}, fmt.Errorf("check username uniqueness: %w", err)
	}

	userRole, valid := vo.ParseUserRole(dto.Role)
	if !valid {
		return User{}, userErrors.InvalidUserRoleError
	}

	passwordHash, err := hasher.HashString(dto.Password)
	if err != nil {
		return User{}, fmt.Errorf("hash password: %w", err)
	}

	user := User{
		ID:           uuid.New(),
		Role:         userRole,
		Login:        dto.Login,
		PasswordHash: passwordHash,
		Username:     dto.Username,
		Email:        dto.Email,
	}

	return user, nil

}

func (u *User) MatchPassword(
	pass string,
) error {
	_, err := rules.ComparePasswordAndHash(pass, u.PasswordHash)
	if err != nil {
		return fmt.Errorf("compare password and hash: %w", err)
	}

	return nil
}

func (u *User) ChangePassword(
	newPass string,
) error {
	hash, err := rules.ComparePasswordAndHash(newPass, u.PasswordHash)
	if err != nil {
		return fmt.Errorf("compare password and hash: %w", err)
	}

	u.PasswordHash = hash
	return nil
}

func (u *User) ChangeLogin(
	ctx context.UserCtx,
	newLogin string,
) error {
	err := rules.CheckLoginUniqueness(ctx, newLogin)
	if err != nil {
		return fmt.Errorf("check login uniqueness: %w", err)
	}
	u.Login = newLogin

	return nil
}

func (u *User) ChangeUsername(
	ctx context.UserCtx,
	newUsername string,
) error {
	err := rules.CheckUsernameUniqueness(ctx, newUsername)
	if err != nil {
		return fmt.Errorf("check username uniqueness: %w", err)
	}
	u.Username = newUsername

	return nil
}

func (u *User) ChangeEmail(
	ctx context.UserCtx,
	newEmail string,
) error {
	err := rules.CheckEmailUniqueness(ctx, newEmail)
	if err != nil {
		return fmt.Errorf("check email uniqiueness: %w", err)
	}

	u.Email = newEmail

	return nil
}
