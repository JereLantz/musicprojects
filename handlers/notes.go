package handlers

import (
	"database/sql"
	"log"
	"musiikkiProjektit/session"
	"musiikkiProjektit/notes"
	"musiikkiProjektit/views/components"
	"musiikkiProjektit/views/pages"
	"net/http"
)

// HandleServeNotes serves the page for saving and displaying notes.
//
// Is a http handler function
func HandleServeNotes(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie(session.SessionTokenName)
	if err != nil{
		w.WriteHeader(200)
		pages.Notes(session.Session{}).Render(r.Context(), w)
		return
	}
	sessionData, err := session.GetSession(cookie.Value)
	if err != nil {
		w.WriteHeader(200)
		pages.Notes(session.Session{}).Render(r.Context(), w)
		return
	}
	w.WriteHeader(200)
	pages.Notes(sessionData).Render(r.Context(), w)
}

// HandleGetSavedNotes is an endpoint that fetches all notes that the user
// has taken. Checks that the user is logged in and the session is valid before
// responding.
//
// requires database connection for fetching the data. Is a http handler function
func HandleGetSavedNotes(db *sql.DB, w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie(session.SessionTokenName)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	token := cookie.Value
	sessionData, err := session.GetSession(token)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	if !sessionData.LoggedIn{
		w.WriteHeader(401)
		return
	}

	// TODO: pitäiskö tässä tarkistaa myös onko session outdated?
	userNotes, err := notes.GetUsersNotes(db, sessionData.Username)
	if err != nil {
		log.Printf("error fetching users notes from db: %s\n", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type","text/html")
	components.NoteDisplay(userNotes).Render(r.Context(), w)
}

// HandleNewNoteForm is a simple endpoint that swaps the new note button
// for the new note form
//
// Is a http handler function
func HandleNewNoteForm(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","text/html")
	components.NewNoteForm(notes.Note{}, []string{}).Render(r.Context(), w)
}

// HandleCreateNewNote is an endpoint that handles the functionality of saving
// new note.
//
// Requires database connection. Is a http handler function
func HandleCreateNewNote(db *sql.DB, w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie(session.SessionTokenName)
	if err != nil {
		//TODO: joku parempi virheellinen session token vastaus
		w.WriteHeader(401)
		return
	}

	userSession, err := session.GetSession(cookie.Value)
	if err != nil {
		//TODO: session doesn't exist.
		log.Printf("Error fetchin session creating new note %s\n", err)
		w.WriteHeader(500)
		return
	}

	if !userSession.LoggedIn{
		//TODO: käyttäjä ei ole kirjautunut sisään. Joku parempi virhe?
		w.WriteHeader(401)
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	newNote := notes.Note{
		Title: r.FormValue("newNoteTitle"),
		Note: r.FormValue("newNote"),
	}

	inputErrors, err := newNote.Validate()
	if err != nil {
		w.Header().Set("Content-Type","text/html")
		components.NewNoteForm(newNote, inputErrors).Render(r.Context(), w)
		return
	}

	err = newNote.SaveNewNote(db, userSession.Username)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type","text/html")
	w.WriteHeader(200)
	components.NoteSavedResp(newNote).Render(r.Context(), w)
}
