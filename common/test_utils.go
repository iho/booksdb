package common

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/iho/booksdb/db"
	"github.com/iho/booksdb/models"
	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Statuses = []models.BookStatusType{models.CheckedIn, models.CheckedOut}

func CreateRandomBook() *models.Book {
	book := new(models.Book)
	book.Title = faker.FirstName()
	book.Author = faker.LastName()
	book.Publisher = faker.Word()
	book.Rating = rand.Intn(3-1) + 1
	book.Status = Statuses[rand.Intn(len(Statuses))]

	return book
}

func CreateRandomBookJSON() *bytes.Buffer {
	book := new(models.Book)
	err := faker.FakeData(&book)
	if err != nil {
		panic(err)
	}

	book.Rating = rand.Intn(3-1) + 1
	book.Status = Statuses[rand.Intn(len(Statuses))]

	jsonValue, _ := json.Marshal(book)
	return bytes.NewBuffer(jsonValue)
}

func GetRandomPort() string {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	defer func() {
		err = listener.Close()
		if err != nil {
			panic(err)
		}
	}()
	_, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		panic(err)
	}

	return port
}

type Repo struct {
	Repo db.BookRepository
	Name string
}

type Suite struct {
	suite.Suite
	Repositories []Repo
	Pool         *dockertest.Pool
	Cancel       *context.CancelFunc
	Resource     *dockertest.Resource
	Context      context.Context
}

func (suite *Suite) Setup() {
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

		port := GetRandomPort()

		dockerOptions := &dockertest.RunOptions{Repository: "mongo", Tag: "4.0", PortBindings: map[dc.Port][]dc.PortBinding{
			dc.Port("27017/tcp"): {{HostPort: port}},
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

func (suite *Suite) TeardDown() {
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
