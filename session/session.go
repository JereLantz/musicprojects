package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"musiikkiProjektit/utils"
	"net/http"
	"time"
)

var sessions = map[string]Session{}

const SESSION_TOKEN_NAME = "session_token"

type Session struct {
	LoggedIn bool
	Username string
	Expiry time.Time
}

func GetSession(id string) (Session, error){
	var requestedSession Session

	requestedSession, ok := sessions[id]
	if !ok {
		return requestedSession, errors.New("No session found with id: " + id)
	}
	return requestedSession,nil
}

func DeleteSession(id string) error{
	delete(sessions, id)
	return nil
}

func UpdateSession(updatedSession Session, id string) error{
	//TODO: tätä tarvitaan ehkä vasta myöhemmin. Kun halutaan tallentaa käyttäjän tilasta enemmän tietoa
	return errors.New("Not yet implemented")
}

func (s Session) isSessionExpired() bool{
	return s.Expiry.Before(time.Now())
}

/*
Creates new session in memory from the Credetials stuct. Can be Credentials struct with null values
if the user is not signed in.
Returns the generated session id
*/
func createSession(creds utils.Credentials) (string, error){
	sessionToken, err := generateSessionId(64)
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(1 * time.Hour)

	sessions[sessionToken] = Session{
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

	reqSession, exists := sessions[sessionToken]
	if !exists{
		return false, nil
	}

	if reqSession.isSessionExpired(){
		delete(sessions, sessionToken)
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
	oldSessionData := sessions[token]
	delete(sessions, token)

	oldSessionData.Expiry = newExpiry
	sessions[token] = oldSessionData

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
				Expires: sessions[newToken].Expiry,
				SameSite: http.SameSiteStrictMode,
			})
			r.AddCookie(&http.Cookie{
				Name: SESSION_TOKEN_NAME,
				Value: newToken,
				Expires: sessions[newToken].Expiry,
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
	sessionDetails := sessions[token]
	delete(sessions, token)

	sessionDetails.LoggedIn = true
	sessionDetails.Username = creds.Username

	sessions[token] = sessionDetails
	return nil
}

func CleanupOutdatedSessions(interval time.Duration){
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C{
		var count int
		log.Println("Cleaning outdated sessions...")
		for token, session := range sessions{
			if session.isSessionExpired(){
				delete(sessions, token)
				count ++
			}
		}
		log.Printf("Cleaned %d outdated sessions", count)
	}
}
