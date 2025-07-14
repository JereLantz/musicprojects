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
	components.NoteDisplay(userNotes).Render(r.Context(), w)
}

func HandleNewNoteForm(w http.ResponseWriter, r *http.Request){
	components.NewNoteForm(notes.Note{}, []string{}).Render(r.Context(), w)
}

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
		components.NewNoteForm(newNote, inputErrors).Render(r.Context(), w)
		return
	}

	err = newNote.SaveNewNote(db, userSession.Username)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	//TODO: joku vastaus joka päivittää uuden noten UI:n jos sen tallennus onnistui
	components.NewNoteForm(newNote, []string{}).Render(r.Context(), w)
}
