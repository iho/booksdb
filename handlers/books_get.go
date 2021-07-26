package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iho/booksdb/db"
)

func (app *App) GetBook(c *gin.Context) {
	id := c.Param("id")

	book, err := app.BookRepository.GetBook(c.Request.Context(), db.ID(id))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}
