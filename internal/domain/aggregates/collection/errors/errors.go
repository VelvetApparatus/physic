package errors

var (
	CollectionIDNotEqualError = newCollectionError("collection id not equal; image cannot be assigned to collection twice")
)

type CollectionError struct {
	Msg string
}

func newCollectionError(msg string) *CollectionError { return &CollectionError{msg} }

func (ce *CollectionError) Error() string { return ce.Msg }
