package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(TeapotHandler),
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func TeapotHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusTeapot)
	writer.Write([]byte(fmt.Sprintf("%v %v", http.StatusTeapot, http.StatusText(http.StatusTeapot))))
}
