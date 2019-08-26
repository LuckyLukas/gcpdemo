package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

var googleOauthConfig *oauth2.Config

var sessions = make(map[string]string)


func initAuth() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/oauth2callback",
		ClientID:	 os.Getenv("GCP_CLIENTID"),
		ClientSecret: os.Getenv("GCP_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

func oauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	session := getUUID()
	state := base64.URLEncoding.EncodeToString([]byte(getUUID()))
	sessions[session] = state
	setAnonymousSessionCookie(w, state)
	u := googleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, u, http.StatusFound)
}

func oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Cookie("SESSION")
	oauthState := sessions[session.Value]

	if r.FormValue("state") != oauthState {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	data, err := getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	authenticatedSession := getUUID()
	sessions[authenticatedSession] = string(data)
	sessions[session.Value] = ""
	setAuthenticatedSessionCookie(w, authenticatedSession)
}

func setAnonymousSessionCookie(w http.ResponseWriter, state string) {
	var expiration = time.Now().Add(1 * time.Minute)
	cookie := http.Cookie{Name: "SESSION", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)
}

func setAuthenticatedSessionCookie(w http.ResponseWriter, state string) {
	var expiration = time.Now().Add(30 * time.Minute)
	cookie := http.Cookie{Name: "SESSION", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)
}

func getUUID() string {
	return uuid.Must(uuid.NewUUID()).String()
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
