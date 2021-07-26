package db

import (
	"context"

	"github.com/iho/booksdb/models"
)

type BookRepository interface {
	AddBook(ctx context.Context, book *models.Book) (ID, error)
	GetBook(ctx context.Context, ID ID) (*models.Book, error)
	DeleteBook(ctx context.Context, ID ID) error
	AllBooks(ctx context.Context) ([]*models.Book, error)
	RemoveAllBooks(ctx context.Context) error
	UpdateBook(
		ctx context.Context,
		ID ID,
		updateFn func(book *models.Book) (*models.Book, error),
	) (*models.Book, error)
}
