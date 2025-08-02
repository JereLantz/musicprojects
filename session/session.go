package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var sessions = map[string]Session{}

// SessionTokenName is the name used for the cookie on the client
const SessionTokenName = "session_token"

type Session struct {
	LoggedIn bool
	Username string
	Expiry time.Time
}

// GetSession attempts to get the session data from the id string
//
// return an error if no session exists
func GetSession(id string) (Session, error){
	var requestedSession Session

	requestedSession, ok := sessions[id]
	if !ok {
		return requestedSession, errors.New("No session found with requested id")
	}
	return requestedSession,nil
}

// GetSessionFromRequest attempts to get the session data from a request
//
// returns bool that whether session exists (true = exists), session data as 
// a session struct, and an error
func GetSessionFromRequest(r *http.Request) (bool, Session, error){
	cookie, err := r.Cookie(SessionTokenName)
	if err != nil {
		return false, Session{}, err
	}

	session, ok := sessions[cookie.Value]
	if !ok {
		return false, Session{}, nil
	}

	return true, session, nil
}

func DeleteSession(id string) error{
	_, err := GetSession(id)
	if err != nil {
		return err
	}

	delete(sessions, id)
	return nil
}

func UpdateSession(updatedSession Session, id string) error{
	_, exists := sessions[id]
	if !exists{
		return errors.New("No session found with given id")
	}

	sessions[id] = updatedSession
	return nil
}

func (s Session) isSessionExpired() bool{
	return s.Expiry.Before(time.Now())
}

/*
Creates new session in memory from the Credetials stuct. Can be Credentials struct with null values
if the user is not signed in.
Returns the generated session id
*/
func createSession() (string, error){
	sessionToken, err := generateSessionId(64)
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(1 * time.Hour)

	sessions[sessionToken] = Session{
		LoggedIn: false,
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
	cookie, err := r.Cookie(SessionTokenName)
	if err != nil {
		if err == http.ErrNoCookie{
			return false, nil
		}
		return false, err
	}
	sessionToken := cookie.Value

	reqSession, err := GetSession(sessionToken)
	if err != nil {
		return false, nil
	}

	if reqSession.isSessionExpired(){
		DeleteSession(sessionToken)
		return false, nil
	}

	return true, nil
}

func refreshSession(w http.ResponseWriter, r *http.Request) error{
	cookie, err := r.Cookie(SessionTokenName)
	if err != nil{
		return err
	}
	token := cookie.Value

	newExpiry := time.Now().Add(1 * time.Hour)
	oldSessionData, err := GetSession(token)
	if err != nil {
		return fmt.Errorf("Failed to fetch session to refresh it. %s", err)
	}

	oldSessionData.Expiry = newExpiry
	//TODO: käytä update session api
	sessions[token] = oldSessionData

	http.SetCookie(w, &http.Cookie{
		Name: SessionTokenName,
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
			newToken, err := createSession()
			if err != nil {
				w.WriteHeader(500)
				log.Printf("Failed to create session: %s\n", err)
				return
			}

			// Delete old cookie, if it exists
			http.SetCookie(w, &http.Cookie{
				Name: SessionTokenName,
				Expires: time.Now(),
			})

			http.SetCookie(w, &http.Cookie{
				Name: SessionTokenName,
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
func SessionLogin(r *http.Request, userName string) error{
	cookie, err := r.Cookie(SessionTokenName)
	if err != nil {
		log.Printf("Failed to fetch the request cookie for logging in and modifying the session %s\n", err)
		return err
	}
	token := cookie.Value
	existingData, err := GetSession(token)
	if err != nil {
		return err
	}

	existingData.LoggedIn = true
	existingData.Username = userName

	UpdateSession(existingData, token)
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
		log.Printf("Cleaning finished. Cleaned %d outdated sessions.", count)
	}
}
