package handlers

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(app *App) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.GET("books", app.ListBooks)
		v1.POST("books", app.CreateBook)
		v1.GET("books/:id", app.GetBook)
		v1.PUT("books/:id", app.UpdateBook)
		v1.DELETE("books/:id", app.DeleteBook)
	}
	return r
}
