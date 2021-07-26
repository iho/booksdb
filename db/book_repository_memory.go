package db

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/iho/booksdb/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemoryBookRepository struct {
	Store   map[ID]*models.Book
	StoreRW *sync.RWMutex
}

func NewMemoryBookRepository() MemoryBookRepository {
	return MemoryBookRepository{
		StoreRW: &sync.RWMutex{},
		Store:   make(map[ID]*models.Book),
	}
}

func (repo MemoryBookRepository) AddBook(ctx context.Context, book *models.Book) (ID, error) {
	id := ID(primitive.NewObjectID().Hex())

	repo.StoreRW.Lock()
	defer repo.StoreRW.Unlock()

	repo.Store[id] = book

	return id, nil
}

func (repo MemoryBookRepository) GetBook(ctx context.Context, id ID) (*models.Book, error) {
	repo.StoreRW.RLock()
	defer repo.StoreRW.RUnlock()

	book, ok := repo.Store[id]
	if ok {
		return book, nil
	}

	return nil, errors.New("book not found")
}

func (repo MemoryBookRepository) DeleteBook(ctx context.Context, id ID) error {
	repo.StoreRW.Lock()
	defer repo.StoreRW.Unlock()

	_, ok := repo.Store[id]
	if ok {
		delete(repo.Store, id)
		return nil
	}

	return errors.New("book not found")
}

func (repo MemoryBookRepository) AllBooks(ctx context.Context) ([]*models.Book, error) {
	var books []*models.Book
	repo.StoreRW.RLock()
	defer repo.StoreRW.RUnlock()
	for _, book := range repo.Store {
		books = append(books, book)
	}

	return books, nil
}

func (repo MemoryBookRepository) RemoveAllBooks(ctx context.Context) error {
	repo.StoreRW.Lock()
	defer repo.StoreRW.Unlock()

	repo.Store = make(map[ID]*models.Book)
	return nil
}

func (repo MemoryBookRepository) UpdateBook(
	ctx context.Context,
	bookID ID,
	updateFn func(book *models.Book) (*models.Book, error),
) (*models.Book, error) {
	var updatedBook *models.Book

	book, err := repo.GetBook(ctx, bookID)
	if err != nil {
		return nil, err
	}

	repo.StoreRW.Lock()
	defer repo.StoreRW.Unlock()

	updatedBook, err = updateFn(book)
	if err != nil {
		return nil, fmt.Errorf("failed to update book: %w", err)
	}

	repo.Store[bookID] = updatedBook

	return updatedBook, nil
}
