package main

import (
	"log"
	"net/http"

	poker "github.com/lethang7794/learn-go-with-tests"
)

const dbFileName = "../game.db.json"

func main() {
	store, cleanup, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()
	server, err := poker.NewPlayerServer(store, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = http.ListenAndServe(":5000", server)
	if err != nil {
		log.Fatal(err)
	}
}
