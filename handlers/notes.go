package handlers

import (
	"musiikkiProjektit/views/notes"
	"net/http"
)

func HandleServeNotes(w http.ResponseWriter, r *http.Request){
	notes.NotesPage().Render(r.Context(), w)
}

