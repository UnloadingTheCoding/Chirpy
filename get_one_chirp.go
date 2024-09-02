package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handleOneChirpReq(w http.ResponseWriter, r *http.Request) {

	chirpID := r.PathValue("chirpID")

	if chirpID != "" {
		id, err := strconv.Atoi(chirpID)
		if err != nil {
			log.Printf("unable to find ID: %s %d", err, id)
		}

		chirp, err := cfg.DB.GetChirp(id)

		if err != nil {
			log.Print(err)
			w.WriteHeader(404)
			return
		}

		data, err := json.Marshal(chirp)

		if err != nil {
			log.Print(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)

	}
}
