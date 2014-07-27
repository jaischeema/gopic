package main

import (
	"github.com/shaoshing/train"
	"github.com/zenazn/goji"
	"net/http"
)

var db = SetupDatabase()
var config = ParseConfig()

func main() {

	goji.Get("/api/photos/:page", PhotosHandler)
	goji.Get("/api/photos/:category/:page", PhotoCategoryHandler)
	goji.Get("/api/photo/:id", PhotoHandler)
	goji.Get("/api/rebuild-index", ReindexHandler)

	train.Config.Verbose = true
	train.ConfigureHttpHandler(nil)

	goji.Serve()
}

func PhotoCategoryHandler(w http.ResponseWriter, r *http.Request) {
}
