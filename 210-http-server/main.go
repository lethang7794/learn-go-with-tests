package main

import (
	"log"
	"net/http"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	database, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	store := NewFileSystemPlayerStore(database)
	server := NewPlayerServer(store)
	err = http.ListenAndServe(":5000", server)
	if err != nil {
		log.Fatal(err)
	}
}
