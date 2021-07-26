package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/iho/booksdb/common"
	"github.com/iho/booksdb/db"
	"github.com/iho/booksdb/models"
	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repo struct {
	Repo db.BookRepository
	Name string
}
type BookRepositoryDBTestSuite struct {
	suite.Suite
	Repositories []Repo
	Pool         *dockertest.Pool
	Cancel       *context.CancelFunc
	Resource     *dockertest.Resource
	Context      context.Context
}

func (suite *BookRepositoryDBTestSuite) SetupSuite() {
	t := suite.T()

	suite.Repositories = make([]Repo, 0)
	suite.Repositories = append(suite.Repositories, Repo{
		Repo: db.NewMemoryBookRepository(),
		Name: "inmemory",
	})

	if !testing.Short() {

		pool, err := dockertest.NewPool("")
		if err != nil {
			t.Fatalf("Could not connect to docker: %s", err)
		}

		suite.Pool = pool

		port := common.GetRandomPort()

		dockerOptions := &dockertest.RunOptions{Repository: "mongo", Tag: "4.0", PortBindings: map[dc.Port][]dc.PortBinding{
			dc.Port("27017/tcp"): []dc.PortBinding{{HostPort: port}},
		}}

		resource, err := pool.RunWithOptions(dockerOptions, func(config *dc.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = dc.RestartPolicy{
				Name: "no",
			}
		})
		if err != nil {
			t.Fatalf("Could not start resource: %s", err)
		}

		suite.Resource = resource

		resource.Expire(60) // Tell docker to hard kill the container in 60 seconds

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		suite.Cancel = &cancel
		suite.Context = ctx

		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:"+port))
		if err != nil {
			t.Fatalf("Cann't connect to mongo container: %s", err)
		}

		err = client.Ping(ctx, options.Client().ReadPreference)
		if err != nil {
			t.Fatal(err)
		}

		suite.Repositories = append(suite.Repositories, Repo{
			Repo: db.NewMongoDBBookRepository(client),
			Name: "MongoDB",
		})
	}
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

func (suite *BookRepositoryDBTestSuite) TestRemoveAllBooks() {
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

func (suite *BookRepositoryDBTestSuite) TeardDownTest() {
	if suite.Cancel != nil {
		(*suite.Cancel)()
	}

	if suite.Pool != nil {
		err := suite.Pool.Purge(suite.Resource)
		if err != nil {
			suite.T().Fatal(err)
		}
	}
}

func TestBookRepositoryDBTestSuite(t *testing.T) {
	suite.Run(t, new(BookRepositoryDBTestSuite))
}
