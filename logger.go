package main

import (
	"log"
)

func CheckError(err error, message string) {
	if err != nil {
		log.Fatalf("Error : %v", message)
	}
}
