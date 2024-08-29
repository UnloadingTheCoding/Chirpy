package main

import "net/http"

func main() {

	serverHandler := http.NewServeMux()
	server := &http.Server{
		Handler: serverHandler,
		Addr:    ":8080",
	}
	serverHandler.Handle("/", http.FileServer(http.Dir('.')))
	serverHandler.Handle("/assets/logo.png", http.FileServer(http.Dir('.')))
	server.ListenAndServe()

}
