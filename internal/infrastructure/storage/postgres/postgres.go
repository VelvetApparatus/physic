package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"physk/internal/domain/aggregates/user"
	"physk/internal/infrastructure/storage"
)

type pgImpl struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) storage.Storage {
	return &pgImpl{db: db}
}

func (p *pgImpl) CreateUser(ctx context.Context, u user.User) error {
	var (
		query = `INSERT INTO user.users (id, role, login, password_hash, username, email) VALUES ($1, $2, $3, $4, $5, $6);`
	)

	_, err := p.db.ExecContext(ctx, query, u.ID, u.Role, u.PasswordHash, u.Username, u.Email)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (p *pgImpl) GetUserByLogin(ctx context.Context, login string) (user.User, error) {
	var (
		query = `SELECT id, role, login, password_hash, username, email FROM user.users WHERE login=$1;`
		usr   user.User
	)

	err := p.db.QueryRowContext(ctx, query, login).
		Scan(
			&usr.ID, &usr.Role, &usr.Login,
			&usr.PasswordHash, &usr.Username,
			&usr.Username, &usr.Email,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.User{}, storage.ErrorNotFound
		}
		return user.User{}, fmt.Errorf("scan: %w", err)
	}

	return usr, nil
}

func (p *pgImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (user.User, error) {
	var (
		query = `SELECT id, role, login, password_hash, username, email FROM user.users WHERE id=$1;`
		usr   user.User
	)

	err := p.db.QueryRowContext(ctx, query, userID).
		Scan(
			&usr.ID, &usr.Role, &usr.Login,
			&usr.PasswordHash, &usr.Username,
			&usr.Username, &usr.Email,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.User{}, storage.ErrorNotFound
		}
		return user.User{}, fmt.Errorf("scan: %w", err)
	}

	return usr, nil
}

var _ storage.Storage = (*pgImpl)(nil)
