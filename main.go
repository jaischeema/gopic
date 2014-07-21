package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

var db = SetupDatabase()
var runningIndex bool = false

func main() {
	r := gin.Default()

	r.GET("/pictures", func(c *gin.Context) {
		photos := LoadPhotos(db, 0, 10)
		c.JSON(200, gin.H{"photos": photos, "status": 200})
	})

	r.GET("/rebuild_index", func(c *gin.Context) {
		if runningIndex {
			c.JSON(200, gin.H{"message": "Reindex already running", "status": 200})
		} else {
			runningIndex = true
			go reindex()
			c.JSON(200, gin.H{"message": "Started reindex", "status": 200})
		}
	})

	r.Run(":" + os.Getenv("PORT"))
}

func reindex() {
	RebuildIndex(db, "/Users/jais/Pictures/Data", ".", false, false)
	fmt.Println("Reindexed data")
	runningIndex = false
}
