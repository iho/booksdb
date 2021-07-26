package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/iho/booksdb/common"
	"github.com/iho/booksdb/handlers"
	"github.com/iho/booksdb/models"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

const (
	JSON_HTTP_HEADER = "application/json"
)

type TestServer struct {
	TS   *httptest.Server
	Name string
}

type BookHandlersTestSuite struct {
	common.Suite
	Servers []TestServer
}

func (suite *BookHandlersTestSuite) SetupSuite() {
	suite.Suite.Setup()

	for _, repo := range suite.Repositories {
		repo := repo

		zapLog, _ := zap.NewDevelopment() // ToDo: change that because don't feel right
		log := zapr.NewLogger(zapLog)

		app := handlers.NewApp(repo.Repo, log)

		suite.Servers = append(suite.Servers, TestServer{
			TS:   httptest.NewServer(handlers.SetupRouter(app)),
			Name: repo.Name,
		})
	}
}

func (suite *BookHandlersTestSuite) getBookFromResponse(resp *http.Response) (*models.Book, error) {
	book := new(models.Book)

	defer resp.Body.Close()

	return book, json.NewDecoder(resp.Body).Decode(book)
}

func (suite *BookHandlersTestSuite) getBooksFromResponse(resp *http.Response) ([]*models.Book, error) {
	books := make([]*models.Book, 0)
	defer resp.Body.Close()
	return books, json.NewDecoder(resp.Body).Decode(&books)
}

func (suite *BookHandlersTestSuite) TestAddBook() {
	t := suite.T()
	t.Parallel()

	for _, server := range suite.Servers {
		server := server
		t.Run("AddBook"+server.Name, func(t *testing.T) {
			t.Parallel()
			book := common.CreateRandomBook()
			jsonValue, _ := json.Marshal(book)
			buf := bytes.NewBuffer(jsonValue)

			resp, err := http.Post(server.TS.URL+"/v1/books", JSON_HTTP_HEADER, buf)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			suite.Assert().NoError(err)
			suite.Assert().Equal(resp.StatusCode, http.StatusCreated)

			book, err = suite.getBookFromResponse(resp)
			suite.Assert().NoError(err)
			suite.Assert().NotEmpty(book.Title)
			suite.Assert().NotEmpty(book.ID)
		})
	}
}

func (suite *BookHandlersTestSuite) TestAddBookFailed() {
	t := suite.T()
	t.Parallel()

	for _, server := range suite.Servers {
		server := server
		t.Run("AddBook"+server.Name, func(t *testing.T) {
			t.Parallel()
			buf := bytes.NewBufferString("{'I am broken JSON")

			resp, err := http.Post(server.TS.URL+"/v1/books", JSON_HTTP_HEADER, buf)
			suite.Assert().Equal(resp.StatusCode, http.StatusBadRequest)
			suite.Assert().NoError(err)
			defer resp.Body.Close()
		})
	}
}

func (suite *BookHandlersTestSuite) TestDeleteBook() {
	t := suite.T()
	t.Parallel()

	for _, server := range suite.Servers {
		server := server
		t.Run("DeleteBook"+server.Name, func(t *testing.T) {
			t.Parallel()
			book := common.CreateRandomBook()
			jsonValue, _ := json.Marshal(book)
			buf := bytes.NewBuffer(jsonValue)

			resp, err := http.Post(server.TS.URL+"/v1/books", JSON_HTTP_HEADER, buf)
			if err != nil {
				t.Fatal(err)
			}

			suite.Assert().NoError(err)
			suite.Assert().Equal(resp.StatusCode, http.StatusCreated)

			book, err = suite.getBookFromResponse(resp)
			suite.Assert().NoError(err)
			suite.Assert().NotEmpty(book.Title)
			suite.Assert().NotEmpty(book.ID)

			client := &http.Client{}

			req, err := http.NewRequest("DELETE", server.TS.URL+"/v1/books/"+book.ID.Hex(), nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err = client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			suite.Assert().Equal(resp.StatusCode, http.StatusNoContent)
		})
	}
}

func (suite *BookHandlersTestSuite) TestDeleteBookFailed() {
	t := suite.T()
	t.Parallel()

	for _, server := range suite.Servers {
		server := server
		t.Run("DeleteBookFailed"+server.Name, func(t *testing.T) {
			t.Parallel()
			client := &http.Client{}

			req, err := http.NewRequest("DELETE", server.TS.URL+"/v1/books/wrong_id", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()
			suite.Assert().NotEqual(resp.StatusCode, http.StatusNoContent)
		})
	}
}

func (suite *BookHandlersTestSuite) TestGetBook() {
	t := suite.T()
	t.Parallel()

	for _, server := range suite.Servers {
		server := server
		t.Run("GetBook"+server.Name, func(t *testing.T) {
			t.Parallel()
			book := common.CreateRandomBook()
			jsonValue, _ := json.Marshal(book)
			buf := bytes.NewBuffer(jsonValue)

			resp, err := http.Post(server.TS.URL+"/v1/books", JSON_HTTP_HEADER, buf)
			if err != nil {
				t.Fatal(err)
			}

			suite.Assert().NoError(err)
			suite.Assert().Equal(resp.StatusCode, http.StatusCreated)

			book, err = suite.getBookFromResponse(resp)
			resp.Body.Close()
			suite.Assert().NoError(err)
			suite.Assert().NotEmpty(book.Title)
			suite.Assert().NotEmpty(book.ID)

			resp, err = http.Get(server.TS.URL + "/v1/books/" + book.ID.Hex())
			suite.Assert().NoError(err)
			suite.Assert().Equal(resp.StatusCode, http.StatusOK)
			defer resp.Body.Close()
		})
	}
}

func (suite *BookHandlersTestSuite) TestGetBookFailed() {
	t := suite.T()
	t.Parallel()

	for _, server := range suite.Servers {
		server := server
		t.Run("GetBookFailed"+server.Name, func(t *testing.T) {
			t.Parallel()
			resp, err := http.Get(server.TS.URL + "/v1/books/wrong_id")
			suite.Assert().NoError(err)
			suite.Assert().NotEqual(resp.StatusCode, http.StatusOK)
		})
	}
}

func (suite *BookHandlersTestSuite) TestListBooks() {
	t := suite.T()
	t.Parallel()

	for _, server := range suite.Servers {
		server := server
		t.Run("ListBooks"+server.Name, func(t *testing.T) {
			t.Parallel()
			book := common.CreateRandomBook()
			jsonValue, _ := json.Marshal(book)
			buf := bytes.NewBuffer(jsonValue)

			resp, err := http.Post(server.TS.URL+"/v1/books", JSON_HTTP_HEADER, buf)
			if err != nil {
				t.Fatal(err)
			}

			suite.Assert().NoError(err)
			suite.Assert().Equal(resp.StatusCode, http.StatusCreated)

			book, err = suite.getBookFromResponse(resp)
			suite.Assert().NoError(err)
			suite.Assert().NotEmpty(book.Title)
			suite.Assert().NotEmpty(book.ID)

			resp, err = http.Get(server.TS.URL + "/v1/books/")
			suite.Assert().NoError(err)
			suite.Assert().Equal(resp.StatusCode, http.StatusOK)

			books, err := suite.getBooksFromResponse(resp)
			suite.Assert().NoError(err)
			suite.Assert().NotEmpty(books)
		})
	}
}

func (suite *BookHandlersTestSuite) TestPutBook() {
	t := suite.T()
	t.Parallel()

	for _, server := range suite.Servers {
		server := server
		t.Run("PutBook"+server.Name, func(t *testing.T) {
			t.Parallel()
			book := common.CreateRandomBook()
			jsonValue, _ := json.Marshal(book)
			buf := bytes.NewBuffer(jsonValue)

			resp, err := http.Post(server.TS.URL+"/v1/books", JSON_HTTP_HEADER, buf)
			if err != nil {
				t.Fatal(err)
			}

			suite.Assert().NoError(err)
			suite.Assert().Equal(resp.StatusCode, http.StatusCreated)

			book, err = suite.getBookFromResponse(resp)
			suite.Assert().NoError(err)
			suite.Assert().NotEmpty(book.Title)
			suite.Assert().NotEmpty(book.ID)

			client := &http.Client{}
			newTitle := "new title"
			book.Title = newTitle
			jsonValue, _ = json.Marshal(book)
			buf = bytes.NewBuffer(jsonValue)

			req, err := http.NewRequest("PUT", server.TS.URL+"/v1/books/"+book.ID.Hex(), buf)
			if err != nil {
				t.Fatal(err)
			}
			resp, err = client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			suite.Assert().NoError(err)
			suite.Assert().Equal(resp.StatusCode, http.StatusCreated)
			book, err = suite.getBookFromResponse(resp)
			suite.Assert().NoError(err)
			suite.Assert().Equal(book.Title, newTitle)
		})
	}
}

func (suite *BookHandlersTestSuite) TeardDownTest() {
	suite.Suite.TeardDown()
}

func TestBookHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(BookHandlersTestSuite))
}
