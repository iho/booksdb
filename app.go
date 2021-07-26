package booksdb

import (
	"net/http"

	"github.com/go-logr/logr"
	"github.com/iho/booksdb/db"
)

type App struct {
	Server         *http.Server
	BookRepository db.BookRepository
	Logger         logr.Logger
}

func NewApp(server *http.Server, bookRepository db.BookRepository, log logr.Logger) *App {
	return &App{
		Server:         server,
		BookRepository: bookRepository,
		Logger:         log,
	}
}

func (app *App) Run() {
	err := app.Server.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}
