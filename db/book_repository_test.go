package db_test

import (
	"errors"
	"testing"

	"github.com/iho/booksdb/common"
	"github.com/iho/booksdb/db"
	"github.com/iho/booksdb/models"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRepositoryDBTestSuite struct {
	common.Suite
}

func (suite *BookRepositoryDBTestSuite) SetupSuite() {
	suite.Suite.Setup()
}

func (suite *BookRepositoryDBTestSuite) TestAddBook() {
	t := suite.T()
	t.Parallel()

	for _, repo := range suite.Repositories {
		repo := repo
		t.Run("AddBook"+repo.Name, func(t *testing.T) {
			t.Parallel()
			book := &models.Book{
				Title: "test book",
			}
			id, err := repo.Repo.AddBook(suite.Context, book)
			suite.Assert().NotEmpty(id)
			suite.Assert().NoError(err)
		})
	}
}

func (suite *BookRepositoryDBTestSuite) TestGetBook() {
	t := suite.T()
	t.Parallel()

	for _, repo := range suite.Repositories {
		repo := repo
		t.Run("GetBook"+repo.Name, func(t *testing.T) {
			t.Parallel()
			title := "test book"
			book := &models.Book{
				Title: title,
			}
			id, err := repo.Repo.AddBook(suite.Context, book)
			suite.Assert().NotEmpty(id)
			suite.Assert().NoError(err)

			insertedBook, err := repo.Repo.GetBook(suite.Context, id)

			suite.Assert().NoError(err)
			suite.Assert().Equal(book.Title, insertedBook.Title)
		})
	}
}

func (suite *BookRepositoryDBTestSuite) TestGetBookFailed() {
	t := suite.T()
	t.Parallel()

	for _, repo := range suite.Repositories {
		repo := repo
		t.Run("GetBookFailed"+repo.Name, func(t *testing.T) {
			t.Parallel()
			id := db.ID(primitive.NewObjectID().Hex())
			_, err := repo.Repo.GetBook(suite.Context, id)
			suite.Assert().Error(err)
		})
	}
}

func (suite *BookRepositoryDBTestSuite) TestAllBooks() {
	t := suite.T()
	t.Parallel()

	for _, repo := range suite.Repositories {
		repo := repo
		t.Run("AllBooks"+repo.Name, func(t *testing.T) {
			t.Parallel()
			title := "test book"
			book := &models.Book{
				Title: title,
			}
			id, err := repo.Repo.AddBook(suite.Context, book)
			suite.Assert().NotEmpty(id)
			suite.Assert().NoError(err)

			books, err := repo.Repo.AllBooks(suite.Context)

			suite.Assert().NoError(err)
			suite.Assert().NotEmpty(books)
		})
	}
}

func (suite *BookRepositoryDBTestSuite) TestUpdateBook() {
	t := suite.T()
	t.Parallel()

	for _, repo := range suite.Repositories {
		repo := repo
		t.Run("UpdateBook"+repo.Name, func(t *testing.T) {
			t.Parallel()
			title := "old title"
			newTitle := "new title"

			book := &models.Book{
				Title: title,
			}
			id, err := repo.Repo.AddBook(suite.Context, book)
			suite.Assert().NotEmpty(id)
			suite.Assert().NoError(err)

			_, err = repo.Repo.UpdateBook(suite.Context, id,
				func(book *models.Book) (*models.Book, error) {
					book.Title = newTitle
					return book, nil
				},
			)
			suite.Assert().NoError(err)
			suite.Assert().NotEmpty(id)

			updatedBook, err := repo.Repo.GetBook(suite.Context, id)
			suite.Assert().NoError(err)

			suite.Assert().Equal(updatedBook.Title, newTitle)
		})
	}
}

func (suite *BookRepositoryDBTestSuite) TestUpdateBookFailed() {
	t := suite.T()
	t.Parallel()

	for _, repo := range suite.Repositories {
		repo := repo
		t.Run("UpdateBookFailed"+repo.Name, func(t *testing.T) {
			t.Parallel()
			wrongId := db.ID(primitive.NewObjectID().Hex())

			_, err := repo.Repo.UpdateBook(suite.Context, wrongId,
				func(book *models.Book) (*models.Book, error) {
					book.Title = "hello"
					return book, nil
				},
			)
			suite.Assert().Error(err)
			suite.Assert().NotEmpty(wrongId)
		})
	}
}

func (suite *BookRepositoryDBTestSuite) TestUpdateBookFailedUpdateFn() {
	t := suite.T()
	t.Parallel()

	for _, repo := range suite.Repositories {
		repo := repo
		t.Run("UpdateBookFailedUpdateFn"+repo.Name, func(t *testing.T) {
			t.Parallel()
			book := &models.Book{
				Title: "Hello",
			}
			id, err := repo.Repo.AddBook(suite.Context, book)
			suite.Assert().NotEmpty(id)
			suite.Assert().NoError(err)

			_, err = repo.Repo.UpdateBook(suite.Context, id,
				func(book *models.Book) (*models.Book, error) {
					return nil, errors.New("I am little sad error ;_;")
				},
			)
			suite.Assert().Error(err)
			suite.Assert().NotEmpty(id)
		})
	}
}

func (suite *BookRepositoryDBTestSuite) TestRemoveAllBooks() { //nolint:tparallel
	t := suite.T()

	for _, repo := range suite.Repositories {
		repo := repo
		t.Run("RemoveAllBooks"+repo.Name, func(t *testing.T) {
			t.Parallel()

			err := repo.Repo.RemoveAllBooks(suite.Context)
			suite.Assert().NoError(err)

			books, err := repo.Repo.AllBooks(suite.Context)
			suite.Assert().NoError(err)
			suite.Assert().Equal(len(books), 0)
		})
	}
}

func (suite *BookRepositoryDBTestSuite) TestDeleteBook() {
	t := suite.T()
	t.Parallel()

	for _, repo := range suite.Repositories {
		repo := repo
		t.Run("DeleteBooks"+repo.Name, func(t *testing.T) {
			t.Parallel()

			title := "test book"
			book := &models.Book{
				Title: title,
			}
			id, err := repo.Repo.AddBook(suite.Context, book)
			suite.Assert().NotEmpty(id)
			suite.Assert().NoError(err)

			err = repo.Repo.DeleteBook(suite.Context, id)
			suite.Assert().NoError(err)

			_, err = repo.Repo.GetBook(suite.Context, id)

			suite.Assert().Error(err)
		})
	}
}

func (suite *BookRepositoryDBTestSuite) TestDeleteBookFailed() {
	t := suite.T()
	t.Parallel()

	for _, repo := range suite.Repositories {
		repo := repo
		t.Run("DeleteBooksFailed"+repo.Name, func(t *testing.T) {
			t.Parallel()

			wrongId := db.ID(primitive.NewObjectID().Hex())
			err := repo.Repo.DeleteBook(suite.Context, wrongId)
			suite.Assert().Error(err)
		})
	}
}

func (suite *BookRepositoryDBTestSuite) TeardDownTest() {
	suite.Suite.TeardDown()
}

func TestBookRepositoryDBTestSuite(t *testing.T) {
	suite.Run(t, new(BookRepositoryDBTestSuite))
}
