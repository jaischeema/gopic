package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

func CheckError(err error, message string) {
	if err != nil {
		panic(fmt.Sprintf("Error : %v", message))
	}
}

func SetupDatabase() *gorm.DB {
	db, err := gorm.Open("sqlite3", "/tmp/gorm.db")
	CheckError(err, "Unable to open database")
	db.DB()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.SingularTable(true)
	db.AutoMigrate(Photo{})
	db.LogMode(true)
	// db.Model(media{}).AddUniqueIndex("idx_media_path", "path")
	// db.Model(media{}).AddUniqueIndex("idx_media_hash", "md5hash")
	return &db
}
