package poker

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

// PlayerStore stores score information about players
type PlayerStore interface {
	GetPlayerScore(player string) (int, bool)
	RecordWin(name string)
	GetLeague() League
}

// Player stores a player with score
type Player struct {
	Name  string
	Score int
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

// NewPlayerServer creates a PlayerServer with routing configured
func NewPlayerServer(store PlayerStore) *PlayerServer {
	server := new(PlayerServer)
	server.store = store

	router := http.NewServeMux()

	router.Handle("/players/", http.HandlerFunc(server.PlayersHandler))
	router.Handle("/league", http.HandlerFunc(server.LeagueHandler))
	router.Handle("/game", http.HandlerFunc(server.GameHandler))
	router.Handle("/ws", http.HandlerFunc(server.WebSocketHandler))

	server.Handler = router

	return server
}

const jsonContentType = "application/json"

func (p *PlayerServer) LeagueHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", jsonContentType)
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(p.store.GetLeague())
}

func (p *PlayerServer) PlayersHandler(writer http.ResponseWriter, request *http.Request) {
	player := strings.TrimPrefix(request.URL.Path, "/players/")

	if request.Method == http.MethodGet {
		p.ShowScore(writer, player)
	}
	if request.Method == http.MethodPost {
		p.ProcessWin(writer, player)
	}
}

func (p *PlayerServer) ShowScore(writer http.ResponseWriter, player string) {
	score, ok := p.store.GetPlayerScore(player)
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(writer, score)
}

func (p *PlayerServer) ProcessWin(writer http.ResponseWriter, player string) {
	writer.WriteHeader(http.StatusAccepted)
	p.store.RecordWin(player)
}

func (p *PlayerServer) GameHandler(writer http.ResponseWriter, request *http.Request) {
	files, err := template.ParseFiles("game.html")
	if err != nil {
		http.Error(writer, fmt.Sprintf("could not parse template file: %s", err.Error()), http.StatusInternalServerError)
	}
	files.Execute(writer, nil)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (p *PlayerServer) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	_, message, _ := conn.ReadMessage()
	p.store.RecordWin(string(message))
}
