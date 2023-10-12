package server

import (
	"context"
	"fmt"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func makeNetBetweenPeers(peers []string) {
	for _, peer := range peers {
		peer := peer
		go makeConnection(peer)
	}
}

func makeConnection(peer string) {
	var (
		ctx             = context.TODO()
		newTransactions = input{}
	)

	for {
		c, _, err := websocket.Dial(ctx, fmt.Sprintf("ws://%s/ws", peer), nil)
		if err = wsjson.Read(ctx, c, &newTransactions); err != nil {
			fmt.Printf("error in websocket client get new transactions %v", err)
		}

		transactionManagerGlobal <- newTransactions
	}
}
