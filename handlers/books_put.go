package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iho/booksdb/db"
	"github.com/iho/booksdb/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (app *App) UpdateBook(c *gin.Context) {
	id := c.Param("id")
	book := new(models.Book)

	err := c.BindJSON(&book)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	book, err = app.BookRepository.UpdateBook(c.Request.Context(), db.ID(id), func(oldBook *models.Book) (*models.Book, error) {
		return book, nil
	})
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	bookID, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	book.ID = bookID
	c.IndentedJSON(http.StatusCreated, book)
}
