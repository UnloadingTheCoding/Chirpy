package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/unloadingthecoding/chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
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
		log.Printf("failure to decode create user request: %s", err)
		w.WriteHeader(500)
		return
	}

	userPW, err := bcrypt.GenerateFromPassword([]byte(create.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to generate pw hash: %s", err)
	}

	user, err := cfg.DB.CreateUser(create.Email, string(userPW))

	if err != nil {
		log.Printf("failure to create user: %s", err)
		w.WriteHeader(500)
		return
	}

	data, err := json.Marshal(map[string]interface{}{"id": user.ID, "email": user.Email})
	if err != nil {
		log.Printf("conversion to json failed: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(data)

}

func (cfg *apiConfig) handleUserLogin(w http.ResponseWriter, r *http.Request) {

	user := database.User{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		log.Printf("failure to decode create user request: %s", err)
		w.WriteHeader(500)
		return
	}

	pwcompare, err := cfg.DB.FindUser(user.Email)
	if err != nil {
		log.Printf("Unable to find user: %s", err)
		w.WriteHeader(500)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(pwcompare.Password), []byte(user.Password))
	if err != nil {
		w.WriteHeader(401)
		return
	}

	data, err := json.Marshal(map[string]interface{}{"id": pwcompare.ID, "email": pwcompare.Email})
	if err != nil {
		log.Printf("conversion to json failed: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
