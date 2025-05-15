package collection

import (
	"github.com/google/uuid"
	collectionDto "physk/internal/domain/aggregates/collection/dto"
	"physk/internal/domain/aggregates/collection/entities"
	collectionErrors "physk/internal/domain/aggregates/collection/errors"
	"time"
)

type Collection struct {
	ID        uuid.UUID
	Name      string
	Images    []entities.Image
	CreatedAt time.Time
}

func NewCollection(dto collectionDto.CreateCollection) Collection {
	return Collection{
		ID:        uuid.New(),
		Name:      dto.Name,
		CreatedAt: time.Now(),
	}
}

func (c *Collection) AddImage(img entities.Image) error {
	if img.CollectionID != c.ID {
		return collectionErrors.CollectionIDNotEqualError
	}
	c.Images = append(c.Images, img)

	return nil
}
