package handlers

import (
	"database/sql"
	"musiikkiProjektit/views/notes"
	"net/http"
)

type Note struct {
	Id int
	Title string
	Note string
}

func HandleServeNotes(w http.ResponseWriter, r *http.Request){
	notes.NotesPage(true).Render(r.Context(), w)
}

func HandleGetSavedNotes(db *sql.DB, w http.ResponseWriter, r *http.Request){
	//TODO:
}
