package entities

import (
	"fmt"
	"github.com/google/uuid"
)

type Image struct {
	ID           uuid.UUID
	CollectionID uuid.UUID
	Name         string
	ContentType  string
	ImageData    []byte
}

func ImageName(collectionID uuid.UUID, imgID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", collectionID, imgID)
}

func (i *Image) ImageName() string { return fmt.Sprintf("%s:%s", i.CollectionID, i.ID) }
