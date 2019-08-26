package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func info(w http.ResponseWriter, r *http.Request) {
	clientid := os.Getenv("GCP_CLIENTID")
	secret := os.Getenv("GCP_SECRET")
	fmt.Fprint(w, "id: %d, s: %d", len(clientid), len(secret))
}

func New() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/",  http.FileServer(http.Dir("templates/")))
	mux.HandleFunc("/info", info)
	// google login
	mux.HandleFunc("/login", oauthGoogleLogin)
	mux.HandleFunc("/oauth2callback", oauthGoogleCallback)

	return mux
}

func main() {
	initAuth()
	server := &http.Server{
		Addr: fmt.Sprintf(":8080"),
		Handler: New(),
	}

	log.Printf("Starting Server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("%v", err)
	} else {
		log.Println("Server closed!")
	}
}
