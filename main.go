package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

var db = SetupDatabase()

func main() {

	r := gin.Default()

	v1 := r.Group("/api")
	{
		v1.GET("/photos", PageParam(), PhotosHandler)
		v1.GET("/photos/*category", PageParam(), PhotoCategoryHandler)
		v1.GET("/photo/:id", PhotoHandler)
		v1.GET("/rebuild-index", ReindexHandler)
	}

	r.GET("/", HomeHandler)
	// r.GET("/login", LoginHandler)

	r.Run(":" + os.Getenv("PORT"))
}
