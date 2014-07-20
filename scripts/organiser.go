package main

import (
	. "github.com/jaischeema/gopic/media"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
)

var db = setupDatabase()

func FindOrCreateMedia(path string) Media {
	var media Media
	db.Where(Media{Path: path}).FirstOrInit(&media)
	if db.NewRecord(media) {
		media.RefreshAttributes()
		db.Save(&media)
	}
	return media
}

func setupDatabase() gorm.DB {
	db, err := gorm.Open("sqlite3", "/tmp/gorm.db")
	// db, err := gorm.Open("postgres", "user=gorm dbname=gorm sslmode=disable")
	checkError(err, "Unable to open database")
	db.DB()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.SingularTable(true)
	db.AutoMigrate(Media{})
	db.LogMode(true)
	// db.Model(media{}).AddUniqueIndex("idx_media_path", "path")
	// db.Model(media{}).AddUniqueIndex("idx_media_hash", "md5hash")
	return db
}

func checkError(err error, message string) {
	if err != nil {
		log.Fatalf("Error : %v", message)
	}
}

func SortMedia(sourceFolder string, destinationFolder string, moveFiles bool, storeInDB bool, findDuplicated bool, generateThumbnails bool) {
	walkFunction := func(path string, info os.FileInfo, err error) error {
		checkError(err, "Unable to transverse source")
		if !info.IsDir() {
			if IsValidMedia(path) {
				path, err := filepath.Abs(path)
				checkError(err, "Invalid file path")
				pic := FindOrCreateMedia(path)
				if moveFiles {
					pic.MoveToDestination(destinationFolder)
					db.Save(&pic)
				}
			}
		}
		return nil
	}
	filepath.Walk(sourceFolder, walkFunction)
}

func main() {
	SortMedia("/source", "./public/system/originals", false, false, false, false)
}
