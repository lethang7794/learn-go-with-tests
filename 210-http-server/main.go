package main

import (
	"log"
	"net/http"
)

func main() {
	store := NewInMemoryPlayerStore()
	handler := &PlayerHandler{store}
	err := http.ListenAndServe(":5000", handler)
	if err != nil {
		log.Fatal(err)
	}
}
