package server

import (
	"context"
	"fmt"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func makeNetBetweenPeers(peers []string) {
	for _, peer := range peers {
		//peer := peer
		go makeConnection(peer)
	}
}

func makeConnection(peer string) {
	var (
		newTransactions = input{}
	)

	//for {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, fmt.Sprintf("ws://%s/ws", peer), nil)
	if err != nil {
		fmt.Println("dial error client websocket", err)
	}

	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	if err = wsjson.Read(ctx, c, newTransactions); err != nil {
		fmt.Printf("error in websocket client get new transactions %v", err)
	}

	transactionManagerGlobal <- newTransactions
	//}
}
