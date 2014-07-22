package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
	"strconv"
	"strings"
)

const perPage = 10

var categoryRegex = regexp.MustCompile(`^(\d{4})(?:\/(\d{1,2})(?:\/(\d{1,2}))?)?$`)

/*************************************** Middleware ******************************/
func PageParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := c.Request.ParseForm()
		CheckError(err, "Error occured while parsing form")
		var page int
		pageParams := c.Request.Form["page"]
		if len(pageParams) > 0 {
			var err error
			page, err = strconv.Atoi(pageParams[0])
			if err != nil {
				page = 0
			}
		} else {
			page = 0
		}
		c.Set("page", page)
		c.Next()
	}
}

/*************************************** Handlers *******************************/
var PhotosHandler = func(c *gin.Context) {
	page := c.MustGet("page").(int)
	photos := LoadPhotos(db, page*perPage, perPage)
	c.JSON(200, gin.H{"photos": photos, "status": 200})
}

var PhotoHandler = func(c *gin.Context) {
	idString := c.Params.ByName("id")
	id, err := strconv.Atoi(idString)
	CheckError(err, "Invalid ID")
	var photo Photo
	db.First(&photo, id)
	c.JSON(200, gin.H{"photo": photo, "status": 200})
}

var HomeHandler = func(c *gin.Context) {
	c.JSON(200, gin.H{"message": "homepage", "status": 200})
}

var PhotoCategoryHandler = func(c *gin.Context) {
	page := c.MustGet("page").(int)
	category := c.Params.ByName("category")
	category = strings.TrimLeft(category, "/")
	category = strings.TrimRight(category, "/")
	matches := categoryRegex.FindAllStringSubmatch(category, -1)
	if matches != nil {
		fmt.Println(matches)
		photos := LoadPhotos(db, page*perPage, perPage)
		c.JSON(200, gin.H{"photos": photos, "status": 200})
	} else {
		c.JSON(400, gin.H{"message": "Invalid category", "status": 400})
	}
}

var runningIndex bool = false
var ReindexHandler = func(c *gin.Context) {
	if runningIndex {
		c.JSON(200, gin.H{"message": "Reindex already running", "status": 200})
	} else {
		runningIndex = true
		go reindexWithCompletion()
		c.JSON(200, gin.H{"message": "Started reindex", "status": 200})
	}
}

func reindexWithCompletion() {
	RebuildIndex(db, "/Users/jais/Pictures/Data", ".", false, false)
	runningIndex = false
}
