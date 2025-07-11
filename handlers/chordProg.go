package handlers

import (
	"log"
	"musiikkiProjektit/session"
	"musiikkiProjektit/views/pages"
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
		w.WriteHeader(500)
		pages.ChordProg(session.Session{}).Render(r.Context(), w)
		return
	}
	pages.ChordProg(sessionData).Render(r.Context(), w)
}
