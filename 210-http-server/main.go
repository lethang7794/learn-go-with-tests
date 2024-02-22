package main

import (
	"log"
	"net/http"
)

func main() {
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	err := http.ListenAndServe(":5000", server)
	if err != nil {
		log.Fatal(err)
	}
}
