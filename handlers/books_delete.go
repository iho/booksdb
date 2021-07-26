package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iho/booksdb/db"
)

func (app *App) DeleteBook(c *gin.Context) {
	id := c.Param("id")

	err := app.BookRepository.DeleteBook(c.Request.Context(), db.ID(id))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusNoContent, gin.H{"message": "book successfully removed"})
}
