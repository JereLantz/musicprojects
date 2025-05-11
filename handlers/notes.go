package handlers

import (
	"database/sql"
	"log"
	"musiikkiProjektit/views/notes"
	"net/http"
)

type Note struct {
	Id int
	Title string
	Note string
}

func HandleServeNotes(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie(SESSION_TOKEN_NAME)
	if err != nil{
		log.Printf("Failed to fetch the session for displaying the notes page. %s\n", err)
		w.WriteHeader(500)
		return
	}
	sessionData := Sessions[cookie.Value]
	notes.NotesPage(sessionData).Render(r.Context(), w)
}

func HandleGetSavedNotes(db *sql.DB, w http.ResponseWriter, r *http.Request){
	//TODO:
}
