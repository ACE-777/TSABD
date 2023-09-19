package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type input struct {
	Body string `json:"body"`
}

func Replace(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(fmt.Sprintf("method must be PUT")))
		return
	}

	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()

	var newBody input
	if err := decode.Decode(&newBody); err != nil {
		w.WriteHeader(http.StatusBadGateway)
	}

	if err := os.WriteFile("../../internal/server/input_body.txt", []byte(newBody.Body), 0777); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err := w.Write([]byte(fmt.Sprintf("Successfully save body")))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Get(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/plain")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(fmt.Sprintf("method must be GET")))
		return
	}

	file, err := os.ReadFile("../../internal/server/input_body.txt")
	if err != nil {
		return
	}

	_, err = w.Write(file)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func StartServer() {
	http.HandleFunc("/replace", Replace)
	http.HandleFunc("/get", Get)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		return
	}
}
