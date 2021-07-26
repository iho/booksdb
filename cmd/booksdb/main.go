package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/iho/booksdb"
	"github.com/iho/booksdb/db"
	"github.com/iho/booksdb/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func main() {
	var log logr.Logger
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}
	config := booksdb.GetConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoDBURL))
	if err != nil {
		panic(fmt.Errorf("can't connect to database: %w", err))
	}

	log = zapr.NewLogger(zapLog)
	app := handlers.NewApp(
		db.NewMongoDBBookRepository(client),
		log,
	)

	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", config.Port),
		Handler:      handlers.SetupRouter(app),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}
