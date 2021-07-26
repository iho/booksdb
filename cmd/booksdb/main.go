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

	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", config.Port),
		Handler:      booksdb.SetupRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log = zapr.NewLogger(zapLog)
	app := booksdb.NewApp(
		server,
		db.NewMongoDBBookRepository(client),
		log,
	)
	app.Run()
}
