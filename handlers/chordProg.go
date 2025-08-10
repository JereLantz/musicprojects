package handlers

import (
	"log"
	"musiikkiProjektit/chorprog"
	"musiikkiProjektit/session"
	"musiikkiProjektit/views/pages"
	"net/http"
)

// HandleServeChordProg renders the Chord Progression lookup page
//
// Requires http response writer and a request pointer
func HandleServeChordProg(w http.ResponseWriter, r *http.Request){
	// TODO: tähän tarkistus eka että onko joku query
	// Käytä getSessionFromRequest?
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
	params := r.URL.Query()

	if len(params) == 0 {
		pages.ChordProg(sessionData).Render(r.Context(), w)
		return
	}

	prog, err := chorprog.GetProgFromParams(params)

	if err != nil {
		w.WriteHeader(400)
		log.Println("HandleServeChordProg(). Fetching chord progression from params:", err)
		return
	}

	pages.ChordProgDisplay(sessionData, prog).Render(r.Context(), w)
}
