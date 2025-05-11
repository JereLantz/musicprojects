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

const SESSION_TOKEN_NAME = "session_token"

var Sessions = map[string]utils.Session{}

/*
Creates new session in memory from the Credetials stuct.
Returns the generated session id
*/
func createSession(creds utils.Credentials) (string, error){
	sessionToken, err := generateSessionId(64)
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(1 * time.Hour)

	Sessions[sessionToken] = utils.Session{
		LoggedIn: false,
		Username: creds.Username,
		Expiry: expiresAt,
	}

	return sessionToken, nil
}

/*
Generates random string the length of len
*/
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
	cookie, err := r.Cookie(SESSION_TOKEN_NAME)
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

	if reqSession.IsSessionExpired(){
		delete(Sessions, sessionToken)
		return false, nil
	}

	return true, nil
}

func refreshSession(w http.ResponseWriter, r *http.Request) error{
	cookie, err := r.Cookie(SESSION_TOKEN_NAME)
	if err != nil{
		return err
	}
	token := cookie.Value

	newExpiry := time.Now().Add(1 * time.Hour)
	oldSessionData := Sessions[token]
	delete(Sessions, token)

	oldSessionData.Expiry = newExpiry
	Sessions[token] = oldSessionData

	http.SetCookie(w, &http.Cookie{
		Name: SESSION_TOKEN_NAME,
		Value: token,
		Expires: newExpiry,
		SameSite: http.SameSiteStrictMode,
	})

	return nil
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
			err = refreshSession(w, r)
			if err != nil {
				log.Printf("Failed to fetch token for refreshing time: %s\n", err)
				return
			}
		} else{
			newToken, err := createSession(utils.Credentials{})
			if err != nil {
				w.WriteHeader(500)
				log.Printf("Failed to create session: %s\n", err)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name: SESSION_TOKEN_NAME,
				Value: newToken,
				Expires: Sessions[newToken].Expiry,
				SameSite: http.SameSiteStrictMode,
			})
		}

		f(w,r)
	}
}

/*
Modifies the session in memory to have the user logged id.
*/
func SessionLogin(r *http.Request, creds utils.Credentials) error{
	cookie, err := r.Cookie(SESSION_TOKEN_NAME)
	if err != nil {
		log.Printf("Failed to fetch the request cookie for logging in and modifying the session %s\n", err)
		return err
	}
	token := cookie.Value
	sessionDetails := Sessions[token]
	delete(Sessions, token)

	sessionDetails.LoggedIn = true
	sessionDetails.Username = creds.Username

	Sessions[token] = sessionDetails
	return nil
}
