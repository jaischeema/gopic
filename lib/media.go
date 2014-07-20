package media

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/hawx/img/exif"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "2006:01:02 15:04:05-07:00"

var ValidExts = []string{".jpg", ".png", ".tiff", ".avi"}

type Media struct {
	Id               int64
	Name             string
	Path             string
	Height           int
	Width            int
	ModificationTime time.Time
	Md5hash          string
}

func (p *Media) SortedPath(root string) string {
	year, month, day := p.ModificationTime.Date()
	path, err := filepath.Abs(fmt.Sprintf("%s/%d/%d/%d/%s", root, year, month, day, p.Name))
	checkError(err, "Invalid path")
	return path
}

func (p *Media) MoveToDestination(destination string) {
	newPath := p.SortedPath(destination)
	err := os.MkdirAll(filepath.Dir(newPath), 0775)
	checkError(err, "Unable to make folder: "+newPath)
	err = os.Rename(p.Path, newPath)
	checkError(err, "Unable to move file: "+p.Path)
	p.Path = newPath
}

func (media *Media) RefreshAttributes() {
	var err error
	path := media.Path
	data := exif.Load(path)
	media.Name = filepath.Base(path)
	media.Height, err = strconv.Atoi(data.Get("ImageHeight"))
	checkError(err, "Invalid image height: "+path)
	media.Width, err = strconv.Atoi(data.Get("ImageWidth"))
	checkError(err, "Invalid image width: "+path)
	media.ModificationTime, err = time.Parse(DateFormat, data.Get("FileModifyDate"))
	checkError(err, "Invalid image date: "+path)
	media.Md5hash = Md5OfFile(path)
}

func checkError(err error, message string) {
	if err != nil {
		log.Fatalf("Error : %v", message)
	}
}

func IsValidMedia(path string) bool {
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
