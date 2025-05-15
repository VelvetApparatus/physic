package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"physk/internal/domain/aggregates/collection"
	"physk/internal/domain/aggregates/collection/entities"
	"physk/internal/domain/aggregates/user"
	"physk/internal/infrastructure/db/storage"
)

type pgImpl struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) storage.Storage {
	return &pgImpl{db: db}
}

func (p *pgImpl) CreateCollection(ctx context.Context, c collection.Collection) error {
	var (
		query = `INSERT INTO collection.collections (id, name, created_at) VALUES ($1, $2, $3);`
	)

	_, err := p.db.ExecContext(ctx, query, c.ID, c.Name, c.CreatedAt)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (p *pgImpl) AddImageToCollection(ctx context.Context, i entities.Image) error {
	var (
		query = `INSERT INTO collection.images (id, colllection_id, name) VALUES ($1, $2, $3)`
	)

	_, err := p.db.ExecContext(ctx, query, i.ID, i.CollectionID, i.Name)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (p *pgImpl) GetCollectionByID(ctx context.Context, collectionID uuid.UUID) (collection.Collection, error) {
	var (
		query = `
SELECT
    c.id, c.name, c.created_at
    FROM collection.collections c 
JOIN collection.images i ON c.id = i.collection_id
WHERE c.id=$1`
		col collection.Collection
	)

	err := p.db.QueryRowContext(ctx, query, collectionID).Scan(&col.ID, &col.Name, &col.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return collection.Collection{}, storage.ErrorNotFound
		}
		return collection.Collection{}, fmt.Errorf("scan: %w", err)
	}

	return col, nil
}

func (p *pgImpl) DeleteCollection(ctx context.Context, collectionID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (p *pgImpl) GetImageIDsByCollectionID(ctx context.Context, collectionID uuid.UUID) ([]uuid.UUID, error) {
	var (
		query = `SELECT id from collection.images WHERE collection_id=$1;`
		id    uuid.UUID
		ids   []uuid.UUID
	)

	rows, err := p.db.QueryContext(ctx, query, collectionID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (p *pgImpl) GetImageByID(ctx context.Context, imgID uuid.UUID) (entities.Image, error) {
	var (
		query = `SELECT id, collection_id, name FROM collection.images WHERE id=$1;`
		img   entities.Image
	)

	err := p.db.QueryRowContext(ctx, query, imgID).Scan(&img.ID, &img.CollectionID, &img.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Image{}, storage.ErrorNotFound
		}
		return entities.Image{}, fmt.Errorf("scan: %w", err)
	}

	return img, nil
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
