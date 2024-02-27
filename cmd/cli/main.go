package main

import (
	"fmt"
	"log"
	"os"

	poker "github.com/lethang7794/learn-go-with-tests"
)

const dbFileName = "../game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println(`Type "{Name} wins" to record a win`)

	store, cleanup, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	pokerHoldem := poker.NewTexasHoldem(store, poker.BlindAlerterFunc(poker.Alerter))
	game := poker.NewCLI(os.Stdin, os.Stdout, pokerHoldem)
	game.PlayPoker()
}
