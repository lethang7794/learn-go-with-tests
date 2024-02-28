package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.HandlerFunc(TeapotHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func TeapotHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusTeapot)
}
