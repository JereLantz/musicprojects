package handlers

import (
	"log"
	"musiikkiProjektit/session"
	"musiikkiProjektit/views/chordProgress"
	"net/http"
)


func HandleServeChordProg(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie(session.SESSION_TOKEN_NAME)
	if err != nil{
		log.Printf("Failed to fetch the session for displaying the chord prog page. %s\n", err)
		w.WriteHeader(500)
		return
	}
	sessionData := session.Sessions[cookie.Value]
	chordprogress.ChordProgPage(sessionData).Render(r.Context(), w)
}
