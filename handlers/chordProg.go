package handlers

import (
	"musiikkiProjektit/session"
	"musiikkiProjektit/views/pages"
	"net/http"
)


// HandleServeChordProg renders the Chord Progression training page to the request
//
// Requires http response writer and a request pointer
func HandleServeChordProg(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie(session.SessionTokenName)
	if err != nil{
		w.WriteHeader(200)
		pages.ChordProg(session.Session{}).Render(r.Context(), w)
		return
	}
	sessionData, err := session.GetSession(cookie.Value)
	if err != nil {
		w.WriteHeader(200)
		pages.ChordProg(session.Session{}).Render(r.Context(), w)
		return
	}
	pages.ChordProg(sessionData).Render(r.Context(), w)
}
