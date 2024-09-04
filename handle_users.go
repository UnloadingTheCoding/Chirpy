package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	user := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		log.Printf("failure to decode create user request: %s", err)
		w.WriteHeader(500)
		return
	}

	defaultExpiration := 60 * 60 * 24
	expiration := time.Duration(defaultExpiration)

	if user.ExpiresInSeconds < defaultExpiration && user.ExpiresInSeconds > 0 {
		expiration = time.Duration(user.ExpiresInSeconds)
	}

	pwcompare, err := cfg.DB.FindUser(user.Email)
	if err != nil {
		log.Printf("Unable to find user: %s", err)
		w.WriteHeader(500)
		return
	}

	id := strconv.Itoa(pwcompare.ID)

	err = bcrypt.CompareHashAndPassword([]byte(pwcompare.Password), []byte(user.Password))
	if err != nil {
		w.WriteHeader(401)
		return
	}

	data, err := json.Marshal(map[string]interface{}{"id": pwcompare.ID, "email": pwcompare.Email, "token": cfg.generateToken(id, expiration*time.Second)})
	if err != nil {
		log.Printf("conversion to json failed: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (cfg *apiConfig) handleUserUpdate(w http.ResponseWriter, r *http.Request) {

	claim := jwt.RegisteredClaims{}

	//strips bearer prefix
	authHeader := r.Header.Get("Authorization")[7:]

	token, err := jwt.ParseWithClaims(authHeader, &claim, func(token *jwt.Token) (interface{}, error) { return []byte(cfg.JWT), nil })
	if err != nil {
		log.Printf("error: %s", err)
		w.WriteHeader(401)
		return
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		w.WriteHeader(401)
		return
	}
	if issuer != string("chirpy") {
		log.Print("invalid issuer")
		w.WriteHeader(401)
		return
	}

	userID, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("error: %s", err)
	}

	user := database.User{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&user)
	if err != nil {
		log.Printf("failure to decode create user request: %s", err)
		w.WriteHeader(500)
		return
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		log.Printf("error: %s", err)
	}

	userPW, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to generate pw hash: %s", err)
	}

	err = cfg.DB.UpdateUser(id, user.Email, string(userPW))
	if err != nil {
		log.Printf("failed to update user: %s", err)
	}

	data, err := json.Marshal(map[string]interface{}{"email": user.Email, "id": id})
	if err != nil {
		log.Print(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func (cfg *apiConfig) generateToken(id string, expiration time.Duration) string {

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiration)),
		Subject:   id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(cfg.JWT))
	if err != nil {
		log.Printf("jwt generation error: %s", err)
	}

	return tokenStr
}
