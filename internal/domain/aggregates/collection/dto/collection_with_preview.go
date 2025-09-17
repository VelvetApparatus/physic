package dto

import "github.com/google/uuid"

type GetCollectionWithPreview struct {
	Items []CollectionWithPreviewItem
}

type CollectionWithPreviewItem struct {
	CollectionID uuid.UUID
	PreviewID    uuid.UUID
	Name         string
}
