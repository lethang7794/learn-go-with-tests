package main

import (
	"fmt"
	"log"
	"os"

	poker "github.com/lethang7794/learn-go-with-tests"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	store, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}

	game := poker.NewCLI(store, os.Stdin)
	game.PlayPoker()
}
