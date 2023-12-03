package server

import (
	"encoding/json"
	"net/http"

	websocket "github.com/gorilla/websocket"
)

type Player struct {
	UID        string `json:"uid"`
	X          int    `json:"x"`
	Y          int    `json:"y"`
	Connection *websocket.Conn
}

type Game struct {
	up      websocket.Upgrader
	players map[string]Player
}

func NewGame() *Game {
	return &Game{
		players: make(map[string]Player),
		up: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (g *Game) broadcastPlayers(player Player) {
	for key, value := range g.players {
		if key == player.UID {
			continue
		}

		payload, err := json.Marshal(value)
		if err != nil {
			continue
		}

		if err = player.Connection.WriteMessage(websocket.TextMessage, payload); err != nil {
			delete(g.players, key)
		}
	}
}

func (g *Game) broadcastPlayer(player Player) {
	payload, err := json.Marshal(player)
	if err != nil {
		return
	}

	for key, value := range g.players {
		if key == player.UID {
			continue
		}

		if err = value.Connection.WriteMessage(websocket.TextMessage, payload); err != nil {
			delete(g.players, key)
		}
	}
}
