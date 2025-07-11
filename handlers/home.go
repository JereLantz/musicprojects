package handlers

import (
	"log"
	"musiikkiProjektit/session"
	"musiikkiProjektit/views/pages"
	"net/http"
)

func HandleServeIndex(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		w.WriteHeader(404)
		return
	}
	cookie, err := r.Cookie(session.SessionTokenName)
	if err != nil{
		w.WriteHeader(200)
		pages.Index(session.Session{}).Render(r.Context(), w)
		return
	}
	sessionData, err := session.GetSession(cookie.Value)
	if err != nil {
		pages.Index(session.Session{}).Render(r.Context(), w)
		log.Printf("could not get session information when serving the index %s\n", err)
		return
	}

	w.WriteHeader(200)
	pages.Index(sessionData).Render(r.Context(), w)
}
