package repository

import (
	"context"

	"cartify/catalogue/domain/model"
)

type BookRepository interface {
	ListBooks(ctx context.Context) ([]*model.Book, error)
	GetBook(ctx context.Context, id int) (*model.Book, error)
}
