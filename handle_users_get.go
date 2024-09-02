package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/unloadingthecoding/chirpy/internal/database"
)

func (cfg *apiConfig) handleUserGetOne(w http.ResponseWriter, r *http.Request) {

	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		log.Printf("Unable to convert string to int: %s", err)
	}

	user, err := cfg.DB.GetUser(userID)
	if err != nil {
		log.Printf("unable to find user: %s", err)
		w.WriteHeader(404)
		return
	}

	data, err := json.Marshal(user)

	if err != nil {
		log.Printf("unknown user format: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(data)

}

func (cfg *apiConfig) handleUserCreate(w http.ResponseWriter, r *http.Request) {

	create := database.User{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&create)
	if err != nil {
		log.Printf("failure to decode create user request: %w", err)
		w.WriteHeader(500)
		return
	}

	user, err := cfg.DB.CreateUser(create.Email)

	if err != nil {
		log.Printf("failure to create user: %w", err)
		w.WriteHeader(500)
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Printf("conversion to json failed: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(data)

}
