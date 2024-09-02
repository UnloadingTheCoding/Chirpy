package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/unloadingthecoding/chirpy/internal/database"
)

func main() {

	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *debug {
		os.Remove("database.json")
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		fmt.Errorf("unable to generate db: %w", err)
	}

	serverHandler := http.NewServeMux()
	server := &http.Server{
		Handler: serverHandler,
		Addr:    ":8080",
	}

	apiConf := &apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	serverHandler.Handle("/app/", apiConf.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	serverHandler.HandleFunc("GET /api/healthz", healthzHandler)
	serverHandler.HandleFunc("GET /admin/metrics", apiConf.hitsHandler)
	serverHandler.HandleFunc("/api/reset", apiConf.resetHitsHandler)
	serverHandler.HandleFunc("GET /api/chirps", apiConf.handlerChirpsGet)
	serverHandler.HandleFunc("POST /api/chirps", apiConf.createChirp)
	serverHandler.HandleFunc("GET /api/chirps/{chirpID}", apiConf.handleOneChirpReq)
	serverHandler.HandleFunc("POST /api/users", apiConf.handleUserCreate)
	serverHandler.HandleFunc("GET /api/users/{userID}", apiConf.handleUserGetOne)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Print(err)
	}

}

func healthzHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func (a *apiConfig) hitsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf(`<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>
</html>`, a.fileserverHits)))

}

func (a *apiConfig) resetHitsHandler(w http.ResponseWriter, r *http.Request) {
	a.fileserverHits = 0
}

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}
