package handlers

import (
	"log"
	"musiikkiProjektit/session"
	"musiikkiProjektit/views/pages"
	"net/http"
)

// HandleServeIndex renders the home page for the request. Checks that the 
// request is for the / path, and not else.
//
// requires http response writer and request
func HandleServeIndex(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		w.WriteHeader(404)
		return
	}

	_, sessionData, err := session.GetSessionFromRequest(r)
	if err != nil {
		log.Println("HandleServeIndex() fetching session data:", err)
	}

	w.WriteHeader(200)
	pages.Index(sessionData).Render(r.Context(), w)
}
