package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func (g *Game) handlePlay(w http.ResponseWriter, r *http.Request) {
	newConnection, err := g.up.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed new connection", http.StatusBadRequest)
		return
	}

	defer func() {
		if err = newConnection.Close(); err != nil {
			log.Println("Error closing connection:", err)
		}
	}()

	g.processingReplicationWork(newConnection)
}

func (g *Game) processingReplicationWork(newConnection *websocket.Conn) {
	for {
		var (
			playerReplica Player

			p []byte
		)

		_, p, err := newConnection.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return
			}

			log.Println("Error reading message:", err)
			continue
		}

		if err = json.Unmarshal(p, &playerReplica); err != nil {
			log.Println("Error unmarshalling JSON:", err)
			continue
		}

		_, ok := g.players[playerReplica.UID]
		playerReplica.Connection = newConnection
		g.players[playerReplica.UID] = playerReplica

		if !ok {
			g.broadcastPlayers(playerReplica)
		}

		g.broadcastPlayer(playerReplica)
	}
}

func StartServer() {
	game := NewGame()
	http.HandleFunc("/ws", game.handlePlay)

	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		fmt.Println("error on 8080:", err)
	}
}
