package main

import (
	"fmt"
	"net/http"
	"strings"
)

func GetPlayerHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "players/")
	if player == "Floyd" {
		fmt.Fprint(w, "10")
		return
	}
	if player == "Pepper" {
		fmt.Fprint(w, "20")
		return
	}
}
