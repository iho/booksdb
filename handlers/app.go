package handlers

import (
	"github.com/go-logr/logr"
	"github.com/iho/booksdb/db"
)

type App struct {
	BookRepository db.BookRepository
	Logger         logr.Logger
}

func NewApp(bookRepository db.BookRepository, log logr.Logger) *App {
	return &App{
		BookRepository: bookRepository,
		Logger:         log,
	}
}
