package collection

import (
	"github.com/google/uuid"
	collectionDto "physk/internal/domain/aggregates/collection/dto"
	"physk/internal/domain/aggregates/collection/entities"
	collectionErrors "physk/internal/domain/aggregates/collection/errors"
	"time"
)

type Collection struct {
	ID             uuid.UUID
	Name           string
	Images         []entities.Image
	CreatedAt      time.Time
	ImagePreviewID uuid.UUID
}

func NewCollection(dto collectionDto.CreateCollection) Collection {
	return Collection{
		ID:        uuid.New(),
		Name:      dto.Name,
		CreatedAt: time.Now(),
	}
}

func (c *Collection) AddImage(
	img entities.Image,
	onPreview bool,
) error {
	if img.CollectionID != c.ID {
		return collectionErrors.CollectionIDNotEqualError
	}
	c.Images = append(c.Images, img)

	if onPreview {
		c.ImagePreviewID = img.ID
	}
	return nil
}
