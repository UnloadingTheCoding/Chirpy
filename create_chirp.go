package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/unloadingthecoding/chirpy/internal/database"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {

	var statusCode int
	var body database.Chirp

	decoder := json.NewDecoder(r.Body)
	chirpMessage := Chirp{}

	err := decoder.Decode(&chirpMessage)
	if err != nil {
		log.Printf("error: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(chirpMessage.Body) > 140 {
		body.Body = "Chirp is too long"
		statusCode = 400
	} else {
		statusCode = 201
		chirpMessage.Body = profaneChecker(chirpMessage.Body)
		body, err = cfg.DB.CreateChirp(chirpMessage.Body)
		if err != nil {
			log.Printf("errroooorrr: %s", err)
		}
	}

	res, err := json.Marshal(body)
	if err != nil {
		log.Printf("error marshalling JSON: %s", err)
		w.WriteHeader(500)
		w.Write(res)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(res)

}

func profaneChecker(check string) string {

	profanity := []string{"kerfuffle", "sharbert", "fornax"}

	checkThis := strings.Split(check, " ")
	for _, pWord := range profanity {
		for i, word := range checkThis {
			if strings.EqualFold(word, pWord) {
				checkThis[i] = "****"
			}
		}
	}
	cleaned := strings.Join(checkThis, " ")

	return cleaned

}
