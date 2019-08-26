package main

import (
	"fmt"
	"log"
	"net/http"
)

func New() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/",  http.FileServer(http.Dir("templates/")))

	// google login
	mux.HandleFunc("/login", oauthGoogleLogin)
	mux.HandleFunc("/oauth2callback", oauthGoogleCallback)

	return mux
}

func main() {
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
