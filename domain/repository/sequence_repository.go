package repository

import "context"

type SequenceRepository interface {
	NextBookId(ctx context.Context) (string, error)
}
