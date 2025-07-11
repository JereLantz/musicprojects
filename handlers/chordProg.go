package handlers

import (
	"musiikkiProjektit/session"
	"musiikkiProjektit/views/pages"
	"net/http"
)


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
