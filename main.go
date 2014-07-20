package main

import "fmt"

var db = SetupDatabase()

func main() {
	SortMedia(db, "./", "/sdsd", false, false, false, false)
	fmt.Println("sddsfsdf")
}
