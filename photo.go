package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/hawx/img/exif"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "2006:01:02 15:04:05-07:00"

var ValidExts = []string{".jpg", ".png", ".tiff", ".avi"}

type Photo struct {
	Id               int64     `json:"id"`
	Name             string    `json:"name"`
	Path             string    `json:"-"`
	Height           int       `json:"height"`
	Width            int       `json:"width"`
	ModificationTime time.Time `json:"modification_time"`
	Md5hash          string    `json:"-"`
	Poster           string    `json:"poster,omitempty" sql:"-"`
	Thumbnail        string    `json:"thumbnail,omitempty" sql:"-"`
}

func (p *Photo) SortedPath(root string) string {
	year, month, day := p.ModificationTime.Date()
	path, err := filepath.Abs(fmt.Sprintf("%s/%d/%d/%d/%s", root, year, month, day, p.Name))
	CheckError(err, "Invalid path")
	return path
}

func (p *Photo) MoveToDestination(destination string) {
	newPath := p.SortedPath(destination)
	err := os.MkdirAll(filepath.Dir(newPath), 0775)
	CheckError(err, "Unable to make folder: "+newPath)
	err = os.Rename(p.Path, newPath)
	CheckError(err, "Unable to move file: "+p.Path)
	p.Path = newPath
}

func (photo *Photo) RefreshAttributes() {
	var err error
	path := photo.Path
	data := exif.Load(path)
	photo.Name = filepath.Base(path)
	photo.Height, err = strconv.Atoi(data.Get("ImageHeight"))
	CheckError(err, "Invalid image height: "+path)
	photo.Width, err = strconv.Atoi(data.Get("ImageWidth"))
	CheckError(err, "Invalid image width: "+path)
	photo.ModificationTime, err = time.Parse(DateFormat, data.Get("FileModifyDate"))
	CheckError(err, "Invalid image date: "+path)
	photo.Md5hash = Md5OfFile(path)
}

func IsValidPhoto(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, b := range ValidExts {
		if b == ext {
			return true
		}
	}
	return false
}

func Md5OfFile(fullpath string) string {
	fi, err := os.Open(fullpath)
	if err != nil {
		return ""
	}
	defer fi.Close()

	r := bufio.NewReader(fi)

	buf := make([]byte, 1024)
	md5sum := md5.New()
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return ""
		}
		if n == 0 {
			break
		}
		md5sum.Write(buf[:n])
	}
	md5 := md5sum.Sum(nil)
	return hex.EncodeToString(md5[:])
}
