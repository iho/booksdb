package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) ListBooks(c *gin.Context) {
	books, err := app.BookRepository.AllBooks(c.Request.Context())
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, books)
}
