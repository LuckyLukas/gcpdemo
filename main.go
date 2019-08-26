package main

import (
	"fmt"
	"log"
	"net/http"
)

func info(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Cookie("SESSION")
	info := sessions[session.Value]
	fmt.Fprint(w, info)
}


func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Cookie("SESSION")
	sessions[session.Value] = ""
	http.Redirect(w, r, "/", 302)
}

func New() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/",  http.FileServer(http.Dir("templates/")))
	mux.HandleFunc("/info", info)

	mux.HandleFunc("/logout", logout)
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
