package main

import (
	"encoding/json"
	"fmt"
	"github.com/zenazn/goji/web"
	"math"
	"net/http"
	"regexp"
	"strconv"
	// "strings"
)

const perPage = 10

var categoryRegex = regexp.MustCompile(`^(\d{4})(?:\/(\d{1,2})(?:\/(\d{1,2}))?)?$`)

/*************************************** Handlers *******************************/
func writeJSONResponse(response map[string]interface{}, writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
	responseString, _ := json.Marshal(response)
	fmt.Fprintf(writer, string(responseString))
}

var PhotosHandler = func(context web.C, writer http.ResponseWriter, request *http.Request) {
	var count int
	db.Table("photo").Count(&count)
	page, err := strconv.Atoi(context.URLParams["page"])
	CheckError(err, "Invalid page param")
	photos := LoadPhotos(db, (page-1)*perPage, perPage)
	totalPages := math.Ceil(float64(count) / float64(perPage))
	writeJSONResponse(map[string]interface{}{"photos": photos, "count": count, "page": page, "total_pages": totalPages}, writer)
}

var PhotoHandler = func(context web.C, writer http.ResponseWriter, request *http.Request) {
	photoId, err := strconv.Atoi(context.URLParams["id"])
	CheckError(err, "Invalid ID")
	var photo Photo
	db.First(&photo, photoId)
	photo.Thumbnail = "/images/2-thumbnail.png"
	photo.Poster = "/images/2-poster.png"
	writeJSONResponse(map[string]interface{}{"photo": photo}, writer)
}

// var PhotoCategoryHandler = func(c *gin.Context) {
// 	page := c.MustGet("page").(int)
// 	category := c.Params.ByName("category")
// 	category = strings.TrimLeft(category, "/")
// 	category = strings.TrimRight(category, "/")
// 	matches := categoryRegex.FindAllStringSubmatch(category, -1)
// 	if matches != nil {
// 		fmt.Println(matches)
// 		photos := LoadPhotos(db, page*perPage, perPage)
// 		c.JSON(200, gin.H{"photos": photos, "status": 200})
// 	} else {
// 		c.JSON(400, gin.H{"message": "Invalid category", "status": 400})
// 	}
// }
//
var runningIndex bool = false
var ReindexHandler = func(context web.C, writer http.ResponseWriter, request *http.Request) {
	if runningIndex {
		writeJSONResponse(map[string]interface{}{"message": "Reindex already running"}, writer)
	} else {
		runningIndex = true
		go reindexWithCompletion()
		writeJSONResponse(map[string]interface{}{"message": "Reindex started"}, writer)
	}
}

func reindexWithCompletion() {
	RebuildIndex(db, config.Get("rootPhotoDirectory"), ".", false, false)
	runningIndex = false
}
