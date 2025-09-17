package context

import "context"

type CollectionCtx struct {
	context.Context
}

func NewCollectionCtx(ctx context.Context) *CollectionCtx {
	return &CollectionCtx{Context: ctx}
}
