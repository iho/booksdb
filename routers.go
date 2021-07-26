package booksdb

import (
	"github.com/gin-gonic/gin"
	"github.com/iho/booksdb/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.GET("books", handlers.ListBooks)
		v1.POST("books", handlers.CreateBook)
		v1.GET("books/:id", handlers.GetBook)
		v1.PUT("books/:id", handlers.UpdateBook)
		v1.DELETE("books/:id", handlers.DeleteBook)
	}
	return r
}
