package repository

import "context"

type UserRepository interface {
	UsernameExists(ctx context.Context, username string) (bool, error)
	LoginExists(ctx context.Context, login string) (bool, error)
	EmailExists(ctx context.Context, email string) (bool, error)
}

/// todo: add provider
