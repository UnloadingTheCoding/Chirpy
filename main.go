package main

import (
	"fmt"
	"net/http"
)

func main() {

	serverHandler := http.NewServeMux()
	server := &http.Server{
		Handler: serverHandler,
		Addr:    ":8080",
	}

	apiConf := &apiConfig{}

	serverHandler.Handle("/app/", apiConf.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	serverHandler.HandleFunc("GET /healthz", healthzHandler)
	serverHandler.HandleFunc("GET /metrics", apiConf.hitsHandler)
	serverHandler.HandleFunc("/reset", apiConf.resetHitsHandler)
	err := server.ListenAndServe()
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
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("Hits: %d", a.fileserverHits)))
}

func (a *apiConfig) resetHitsHandler(w http.ResponseWriter, r *http.Request) {
	a.fileserverHits = 0
}

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}
