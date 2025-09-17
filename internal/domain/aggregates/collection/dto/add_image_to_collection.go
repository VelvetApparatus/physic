package dto

import "github.com/google/uuid"

type AddImageToCollection struct {
	CollectionID uuid.UUID
	Name         string
	ContentType  string
	Data         []byte
	IsPreview    bool
}
