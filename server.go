package poker

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

const htmlTemplatePath = "game.html"

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
	template *template.Template
	game     Game
}

// NewPlayerServer creates a PlayerServer with routing configured
func NewPlayerServer(store PlayerStore, game Game) (*PlayerServer, error) {
	server := new(PlayerServer)
	server.store = store
	server.game = game

	tmpl, err := template.ParseFiles(htmlTemplatePath)
	if err != nil {
		return nil, fmt.Errorf("could not parse template file: %s", err)
	}
	server.template = tmpl

	router := http.NewServeMux()

	router.Handle("/players/", http.HandlerFunc(server.PlayersHandler))
	router.Handle("/league", http.HandlerFunc(server.LeagueHandler))
	router.Handle("/game", http.HandlerFunc(server.GameHandler))
	router.Handle("/ws", http.HandlerFunc(server.WebSocketHandler))

	server.Handler = router

	return server, nil
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
	p.template.Execute(writer, nil)
}

type playerServerWS struct {
	conn *websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func newPlayerServerWS(w http.ResponseWriter, r *http.Request) *playerServerWS {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &playerServerWS{conn}
}

func (w playerServerWS) WaitForMessage() string {
	_, message, err := w.conn.ReadMessage()
	if err != nil {
		log.Print(err)
	}
	return string(message)
}

func (w playerServerWS) Write(p []byte) (n int, err error) {
	err = w.conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (p *PlayerServer) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	ws := newPlayerServerWS(w, r)

	message := ws.WaitForMessage()
	numberOfPlayer, _ := strconv.Atoi(message)
	p.game.Start(numberOfPlayer, ws)

	message = ws.WaitForMessage()
	p.game.Finish(message)
}
