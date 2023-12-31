package repository

import (
	"context"
	"fmt"

	"cartify/catalogue/domain/domainerror"
	"cartify/catalogue/domain/repository"

	"github.com/jmoiron/sqlx"
)

const nextSequenceSQL = "UPDATE sequences SET sequence = (LAST_INSERT_ID(sequence) + 1) WHERE name = ?"

type SequenceRepository struct {
	db *sqlx.DB
}

func NewSequenceRepository(db *sqlx.DB) repository.SequenceRepository {
	return &SequenceRepository{db: db}
}

func (s *SequenceRepository) NextBookId(ctx context.Context) (string, error) {
	next, err := s.next("id")
	if err != nil {
		return "", domainerror.NewInternalServerError(err.Error(), err)
	}
	if next > 999999 {
		return "", domainerror.NewInternalServerError("id generation reached to limit", nil)
	}
	return fmt.Sprintf("%07d", next), nil
}

func (s *SequenceRepository) next(name string) (int64, error) {
	result, err := s.db.Exec(nextSequenceSQL, name)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
