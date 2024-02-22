package main

import (
	"log"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(GetPlayerHandler)
	err := http.ListenAndServe(":5000", handler)
	if err != nil {
		log.Fatal(err)
	}
}
