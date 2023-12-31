package book

import (
	"cartify/catalogue/domain/model"
	"cartify/catalogue/domain/repository"
	"context"
)

type ListBooksParams struct {
}

type ListBooks func(ctx context.Context) ([]*model.Book, error)

func NewListBooks(bookRepository repository.BookRepository) ListBooks {
	return func(ctx context.Context) ([]*model.Book, error) {
		return bookRepository.ListBooks(ctx)
	}
}
