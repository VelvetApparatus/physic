package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"physk/internal/domain/aggregates/user/repository"
	vo "physk/internal/domain/aggregates/user/value_objects"
)

type pgImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &pgImpl{db: db}
}

func (p *pgImpl) UsernameExists(ctx context.Context, username vo.UserName) (bool, error) {
	var (
		query  = `SELECT EXISTS(SELECT 1 FROM user.users WHERE username=$1);`
		exists bool
	)

	err := p.db.QueryRowContext(ctx, query, username.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("scan: %w", err)
	}

	return exists, nil
}

func (p *pgImpl) LoginExists(ctx context.Context, login vo.UserLogin) (bool, error) {
	var (
		query  = `SELECT EXISTS(SELECT 1 FROM user.users WHERE login=$1);`
		exists bool
	)

	err := p.db.QueryRowContext(ctx, query, login.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("scan: %w", err)
	}

	return exists, nil
}

func (p *pgImpl) EmailExists(ctx context.Context, email vo.UserEmail) (bool, error) {
	var (
		query  = `SELECT EXISTS(SELECT 1 FROM user.users WHERE email=$1);`
		exists bool
	)

	err := p.db.QueryRowContext(ctx, query, email.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("scan: %w", err)
	}

	return exists, nil
}

var _ repository.UserRepository = (*pgImpl)(nil)
