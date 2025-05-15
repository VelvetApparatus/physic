package user

import (
	"fmt"
	"github.com/google/uuid"
	"physk/internal/domain/aggregates/user/context"
	userDto "physk/internal/domain/aggregates/user/dto"
	userErrors "physk/internal/domain/aggregates/user/errors"
	"physk/internal/domain/aggregates/user/rules"
	vo "physk/internal/domain/aggregates/user/value_objects"
)

type User struct {
	ID           uuid.UUID
	Role         vo.UserRole
	Login        vo.UserLogin
	PasswordHash vo.UserPasswordHash
	Username     vo.UserName
	Email        vo.UserEmail
}

func CreateUser(
	ctx context.UserCtx,
	dto userDto.CreateUserDTO,
) (User, error) {

	login, valid := vo.ParseLogin(dto.Login)
	if !valid {
		return User{}, userErrors.InvalidUserLoginError
	}

	userRole, valid := vo.ParseUserRole(dto.Role)
	if !valid {
		return User{}, userErrors.InvalidUserRoleError
	}

	passwordHash, valid := vo.ParsePassword(dto.Password)
	if !valid {
		return User{}, userErrors.InvalidPasswordError
	}

	email, valid := vo.ParseEmail(dto.Email)
	if !valid {
		return User{}, userErrors.InvalidUserEmailError
	}

	name, valid := vo.ParseUsername(dto.Username)
	if !valid {
		return User{}, userErrors.InvalidUsernameError
	}

	err := rules.CheckLoginUniqueness(ctx, login)
	if err != nil {
		return User{}, fmt.Errorf("check login uniqueness: %w", err)
	}

	err = rules.CheckEmailUniqueness(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("check email uniquenes: %w", err)
	}

	err = rules.CheckUsernameUniqueness(ctx, name)
	if err != nil {
		return User{}, fmt.Errorf("check username uniqueness: %w", err)
	}

	user := User{
		ID:           uuid.New(),
		Role:         userRole,
		Login:        login,
		PasswordHash: passwordHash,
		Username:     name,
		Email:        email,
	}

	return user, nil

}

func (u *User) MatchPassword(
	pass string,
) error {
	_, err := rules.ComparePasswordAndHash(pass, u.PasswordHash.String())
	if err != nil {
		return fmt.Errorf("compare password and hash: %w", err)
	}

	return nil
}

func (u *User) ChangePassword(
	newPass string,
) error {

	hash, valid := vo.ParsePassword(newPass)
	if !valid {
		return userErrors.InvalidPasswordError
	}

	u.PasswordHash = hash
	return nil
}

func (u *User) ChangeLogin(
	ctx context.UserCtx,
	newLogin string,
) error {
	login, valid := vo.ParseLogin(newLogin)
	if !valid {
		return userErrors.InvalidUserLoginError
	}
	err := rules.CheckLoginUniqueness(ctx, login)
	if err != nil {
		return fmt.Errorf("check login uniqueness: %w", err)
	}
	u.Login = login

	return nil
}

func (u *User) ChangeUsername(
	ctx context.UserCtx,
	newUsername string,
) error {
	username, valid := vo.ParseUsername(newUsername)
	if !valid {
		return userErrors.InvalidUsernameError
	}

	err := rules.CheckUsernameUniqueness(ctx, username)
	if err != nil {
		return fmt.Errorf("check username uniqueness: %w", err)
	}
	u.Username = username

	return nil
}

func (u *User) ChangeEmail(
	ctx context.UserCtx,
	newEmail string,
) error {
	email, valid := vo.ParseEmail(newEmail)
	if !valid {
		return userErrors.InvalidUserEmailError
	}

	err := rules.CheckEmailUniqueness(ctx, email)
	if err != nil {
		return fmt.Errorf("check email uniqiueness: %w", err)
	}

	u.Email = email

	return nil
}
