package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"musiikkiProjektit/utils"
	"net/http"
	"time"
)

var Sessions = map[string]Session{}

type Session struct {
	LoggedIn bool
	Username string
	Expiry time.Time
}

func (s Session) isSessionExpired() bool{
	return s.Expiry.Before(time.Now())
}

/*
Creates new session and stores it in memory from the Credetials stuct.
Then returns the generated session id
*/
func createSession(creds utils.Credentials) (string, error){
	sessionToken, err := generateSessionId(64)
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(1 * time.Hour)

	Sessions[sessionToken] = Session{
		LoggedIn: false,
		Username: creds.Username,
		Expiry: expiresAt,
	}

	return sessionToken, nil
}

func generateSessionId(len int) (string, error){
	b := make([]byte, len)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

/*
Check whether the request contains a valid session.
*/
func checkForValidSession(r *http.Request) (bool, error){
	//TODO:
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie{
			return false, nil
		}
		return false, err
	}
	sessionToken := cookie.Value

	reqSession, exists := Sessions[sessionToken]
	if !exists{
		return false, nil
	}

	if reqSession.isSessionExpired(){
		delete(Sessions, sessionToken)
		return false, nil
	}

	refreshSession()

	return true, nil
}

func refreshSession(){
	//TODO:
	log.Println("refresh")
}

func HandleSessionMiddleware(f http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		hasSession, err := checkForValidSession(r)
		if err != nil {
			w.WriteHeader(500)
			log.Printf("Failed to validate request session: %s\n", err)
			return
		}

		if hasSession{
			refreshSession()
		} else{
			token, err := createSession(utils.Credentials{})
			if err != nil {
				w.WriteHeader(500)
				log.Printf("Failed to create session: %s\n", err)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name: "session_token",
				Value: token,
				Expires: Sessions[token].Expiry,
				SameSite: http.SameSiteStrictMode,
			})
		}

		f(w,r)
	}
}
