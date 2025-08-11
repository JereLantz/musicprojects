package handlers

import (
	"database/sql"
	"log"
	"musiikkiProjektit/notes"
	"musiikkiProjektit/session"
	"musiikkiProjektit/views/components"
	"musiikkiProjektit/views/pages"
	"net/http"
	"strconv"
	"time"
)

// HandleServeNotes serves the page for saving and displaying notes.
//
// Is a http handler function
func HandleServeNotes(w http.ResponseWriter, r *http.Request){
	_, sessionData, err := session.GetSessionFromRequest(r)
	if err != nil {
		log.Println("HandleServeNotes() fetching session data:",err)
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

	id, err := newNote.SaveNewNote(db, userSession.Username)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	newNote.Id = id
	newNote.Created = time.Now()

	w.Header().Set("Content-Type","text/html")
	w.WriteHeader(200)
	components.NoteSavedResp(newNote).Render(r.Context(), w)
}


// HandleDeleteNote is the endpoint for deleting a note.
//
// http handler function
func HandleDeleteNote(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	noteIdStr := r.PathValue("id")
	noteId, err := strconv.Atoi(noteIdStr)
	if err != nil {
		log.Println("HandleDeleteNote() parsing noteId from path:", err)
		w.WriteHeader(400)
		return
	}

	cookie, err := r.Cookie(session.SessionTokenName)
	if err != nil {
		log.Println("HandleDeleteNote() getting the requests cookie:", err)
		w.WriteHeader(401)
		return
	}

	session, err := session.GetSession(cookie.Value)
	if err != nil {
		log.Println("HandleDeleteNote() getting the session:", err)
		w.WriteHeader(401)
		return
	}

	err = notes.DeleteNote(db, noteId, session.Username)
	if err != nil {
		log.Println("HandleDeleteNote() calling the DeleteNote():", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

// HandleEditNote serves the edit page
//
// http handler function
func HandleEditNote(db *sql.DB, w http.ResponseWriter, r *http.Request){
	noteIdStr := r.PathValue("id")
	noteId, err := strconv.Atoi(noteIdStr)
	if err != nil {
		log.Println("HandleEditNote() parsing note id from string:", err)
		w.WriteHeader(400)
		return
	}

	cookie, err := r.Cookie(session.SessionTokenName)
	if err != nil {
		log.Println("HandleEditNote() getting cookie:", err)
		w.WriteHeader(400)
		return
	}

	session, err := session.GetSession(cookie.Value)
	if err != nil || !session.LoggedIn {
		log.Println("HandleEditNote() getting session:", err)
		w.WriteHeader(400)
		return
	}

	note, err := notes.GetNote(db, noteId, session.Username)
	if err != nil{
		log.Println("HandleEditNote() getting note:", err)
		w.WriteHeader(400)
		return
	}
	pages.EditNote(session, note).Render(r.Context(), w)
}
