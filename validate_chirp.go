package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func validate_chirp(w http.ResponseWriter, r *http.Request) {

	type chirp struct {
		ChirpMessage string `json:"body"`
		Validation   bool   `json:"valid"`
		ErrorMsg     string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	chirpMessage := chirp{}
	err := decoder.Decode(&chirpMessage)
	if err != nil {
		log.Printf("error: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(chirpMessage.ChirpMessage) > 140 {
		chirpMessage.Validation = false
		chirpMessage.ErrorMsg = "Chirp is too long"
		w.WriteHeader(400)
		w.Write([]byte(chirpMessage.ErrorMsg))
		return
	}

	chirpMessage.Validation = true

	res, err := json.Marshal(chirpMessage)
	if err != nil {
		log.Printf("error marshalling JSON: %s", err)
		w.WriteHeader(500)
		w.Write(res)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(res)

}
