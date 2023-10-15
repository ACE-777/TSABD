package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	snap                     = "{}"
	wal                      []input
	IDTransaction            uint64
	transactionManagerGlobal chan input

	timer *time.Timer

	clock map[string]uint64

	source = "Diagilev"

	peers []string
)

type input struct {
	Source  string `json:"source"`
	Id      uint64 `json:"id"`
	Payload string `json:"payload"`
}

func replace(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(fmt.Sprintf("method must be PUT")))
		return
	}

	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()

	var newTransaction input
	if err := decode.Decode(&newTransaction); err != nil {
		w.WriteHeader(http.StatusBadGateway)
	}

	var wg sync.WaitGroup
	wg.Wait()
	wg.Add(1)
	go func() {
		transactionManagerGlobal <- newTransaction
		wg.Done()
	}()

	if err := os.WriteFile("internal/server/input_body.txt", []byte(newTransaction.Payload), 0777); err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	_, err := w.Write([]byte(fmt.Sprintf("Successfully save body")))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func get(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/plain")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(fmt.Sprintf("method must be GET")))
		return
	}
	//file, err := os.ReadFile("internal/server/input_body.txt")
	//if err != nil {
	//	return
	//}
	//
	//_, err = w.Write(file)
	_, err := w.Write([]byte(snap))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func test(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/plain")

	tmpl, err := template.ParseFiles("static/templates/index.html")
	if err != nil {
		http.Error(w, "Error in parsing index.html", http.StatusBadGateway)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func vClock(w http.ResponseWriter, r *http.Request) {
	out := "\n"
	for key, value := range clock {
		out = out + key + " " + strconv.Itoa(int(value)) + "\n"
	}

	w.Write([]byte(out))
	w.WriteHeader(http.StatusOK)
}

func ws(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		OriginPatterns:     []string{"*"},
	})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	defer func(c *websocket.Conn, code websocket.StatusCode, reason string) {
		if err = c.Close(code, reason); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}(c, websocket.StatusInternalError, "the sky is falling")

	if err = wsjson.Write(r.Context(), c, wal); err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	//ws- это websocket handler, по которому мы отправляем транзакции,
	// а надо еще в отдельной горутине поднять клиента для принятия от всех peer-ов транзакции. То есть для каждого
	//другого пира нужна горутина с клиентом
}

func StartServer() {
	clock = make(map[string]uint64)
	clock[source] = 0
	timer = time.NewTimer(10 * time.Minute)
	transactionManagerGlobal = make(chan input)
	go makeLog(1)
	peers = append(peers, "127.0.0.1:8080", "127.0.0.1:8081")
	makeNetBetweenPeers(peers)

	http.HandleFunc("/replace", replace)
	http.HandleFunc("/get", get)
	http.HandleFunc("/test", test)
	http.HandleFunc("/vclock", vClock)
	http.HandleFunc("/ws", ws)

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println("error on 8080:", err)
		}
	}()

	go func() {
		if err := http.ListenAndServe(":8081", nil); err != nil {
			fmt.Println("error on:", err)
		}
	}()

	time.Sleep(300 * time.Second)
}
