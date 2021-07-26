package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/iho/booksdb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBBookRepository struct {
	Client *mongo.Client
}

func NewMongoDBBookRepository(client *mongo.Client) MongoDBBookRepository {
	return MongoDBBookRepository{
		Client: client,
	}
}

func (repo MongoDBBookRepository) getBookCollection() *mongo.Collection {
	return repo.Client.Database(DatabaseName).Collection(BookCollectionName)
}

func (repo MongoDBBookRepository) AddBook(ctx context.Context, book *models.Book) (ID, error) {
	BookID := ID("")

	result, err := repo.getBookCollection().InsertOne(ctx, book)
	if err != nil {
		return BookID, fmt.Errorf("can't insert a book: %w", err)
	}

	return ID(result.InsertedID.(primitive.ObjectID).Hex()), nil
}

func (repo MongoDBBookRepository) GetBook(ctx context.Context, id ID) (*models.Book, error) {
	book := &models.Book{}

	bookID, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return book, fmt.Errorf("ID is not a valid hex string: %w", err)
	}

	filter := bson.M{"_id": bookID}

	err = repo.getBookCollection().FindOne(ctx, filter).Decode(book)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return book, fmt.Errorf("can't find a book: %w", err)
	} else if err != nil {
		return book, fmt.Errorf("some database error has occurred: %w", err)
	}

	return book, nil
}

func (repo MongoDBBookRepository) AllBooks(ctx context.Context) ([]*models.Book, error) {
	var books []*models.Book

	filter := bson.M{}

	cur, err := repo.getBookCollection().Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("some database error has occurred: %w", err)
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		result := &models.Book{}

		err := cur.Decode(&result)
		if err != nil {
			return nil, fmt.Errorf("can't decode a book: %w", err)
		}
		books = append(books, result)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("some database error has occurred: %w", err)
	}

	return books, nil
}

func (repo MongoDBBookRepository) RemoveAllBooks(ctx context.Context) error {
	err := repo.getBookCollection().Drop(ctx)
	if err != nil {
		return fmt.Errorf("can't drop a book collection: %w", err)
	}

	return nil
}

func (repo MongoDBBookRepository) UpdateBook(
	ctx context.Context,
	bookID ID,
	updateFn func(book *models.Book) (*models.Book, error),
) (*models.Book, error) {
	var updatedBook *models.Book
	err := repo.Client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		book, err := repo.GetBook(ctx, bookID)
		if err != nil {
			sessionContext.AbortTransaction(ctx)

			return err
		}

		updatedBook, err = updateFn(book)
		if err != nil {
			sessionContext.AbortTransaction(ctx)

			return err
		}
		filter := bson.M{"_id": bson.M{"$eq": book.ID}}

		opts := options.Replace().SetUpsert(true)
		_, err = repo.getBookCollection().ReplaceOne(ctx, filter, updatedBook, opts)
		if err != nil {
			return err
		}

		err = sessionContext.CommitTransaction(ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return updatedBook, fmt.Errorf("failed to update book: %w", err)
	}

	return updatedBook, nil
}
