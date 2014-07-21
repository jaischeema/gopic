package main

import (
	"github.com/jinzhu/gorm"
	"os"
	"path/filepath"
)

func FindOrCreatePhoto(db *gorm.DB, path string) Photo {
	var media Photo
	db.Where(Photo{Path: path}).FirstOrInit(&media)
	if db.NewRecord(media) {
		media.RefreshAttributes()
		db.Save(&media)
	}
	return media
}

func RebuildIndex(db *gorm.DB, sourceFolder string, destinationFolder string, moveFiles bool, generateThumbnails bool) {
	walkFunction := func(path string, info os.FileInfo, err error) error {
		CheckError(err, "Unable to transverse source")
		if !info.IsDir() {
			if IsValidPhoto(path) {
				path, err := filepath.Abs(path)
				CheckError(err, "Invalid file path")
				pic := FindOrCreatePhoto(db, path)
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
