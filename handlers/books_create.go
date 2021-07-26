package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iho/booksdb/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (app *App) CreateBook(c *gin.Context) {
	book := new(models.Book)

	err := c.BindJSON(&book)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return
	}

	id, err := app.BookRepository.AddBook(c.Request.Context(), book)
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
