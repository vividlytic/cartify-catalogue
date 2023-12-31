package repository

import (
	"cartify/catalogue/domain/domainerror"
	"cartify/catalogue/domain/model"
	"cartify/catalogue/domain/repository"
	"context"

	"github.com/jmoiron/sqlx"
)

type BookRepository struct {
	db                 *sqlx.DB
	sequenceRepository repository.SequenceRepository
}

func NewBookRepository(db *sqlx.DB, sequenceRepository repository.SequenceRepository) repository.BookRepository {
	return &BookRepository{
		db:                 db,
		sequenceRepository: sequenceRepository,
	}
}

func (b *BookRepository) ListBooks(ctx context.Context) ([]*model.Book, error) {
	books := []*model.Book{}
	query := "SELECT * FROM books"
	err := b.db.Select(&books, query)
	if err != nil {
		return nil, domainerror.NewInternalServerError(err.Error(), err)
	}
	return books, nil
}

func (b *BookRepository) GetBook(ctx context.Context, id int) (*model.Book, error) {
	book := &model.Book{}
	query := "SELECT * FROM books WHERE id=?"
	err := b.db.Get(book, query, id)
	if err != nil {
		return nil, domainerror.NewInternalServerError(err.Error(), err)
	}
	return book, nil
}
