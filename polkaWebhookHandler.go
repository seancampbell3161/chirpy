package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type polkaEvent struct {
	Event string    `json:"event"`
	Data  polkaData `json:"data"`
}

type polkaData struct {
	UserID int `json:"user_id"`
}

func (cfg *apiConfig) polkaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	polka := polkaEvent{}
	err := decoder.Decode(&polka)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
	}

	if polka.Event != "user.upgraded" {
		w.WriteHeader(200)
		return
	}

	user, err := cfg.DB.UpdateUserMembership(polka.Data.UserID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(404)
		return
	}

	data, err := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err = w.Write(data)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
}
