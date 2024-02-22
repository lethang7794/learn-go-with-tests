package main

import (
	"log"
	"net/http"
)

func main() {
	handler := &PlayerHandler{}
	err := http.ListenAndServe(":5000", handler)
	if err != nil {
		log.Fatal(err)
	}
}
