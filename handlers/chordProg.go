package handlers

import (
	"log"
	"musiikkiProjektit/session"
	"musiikkiProjektit/views/chordProgress"
	"net/http"
)


func HandleServeChordProg(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie(session.SessionTokenName)
	if err != nil{
		log.Printf("Failed to fetch the session for displaying the chord prog page. %s\n", err)
		w.WriteHeader(500)
		return
	}
	sessionData, err := session.GetSession(cookie.Value)
	if err != nil {
		//TODO: session doesn't exist.
		log.Printf("Error fetchin session when displaying chord prog page %s\n", err)
		w.WriteHeader(500)
		return
	}
	chordprogress.ChordProgPage(sessionData).Render(r.Context(), w)
}
