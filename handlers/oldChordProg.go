package handlers

import (
	"log"
	"musiikkiProjektit/session"
	"musiikkiProjektit/views/pages"
	"net/http"
)

func HandleServeOldChordProg(w http.ResponseWriter, r *http.Request){
	_, sessionData, err := session.GetSessionFromRequest(r)
	if err != nil {
		log.Println("HandleServeOldChordProg() fetching session data: ", err)
	}

	pages.OldChordProg(sessionData).Render(r.Context(), w)
}
