package handlers

import (
	"database/sql"
	"log"
	"musiikkiProjektit/utils"
	"musiikkiProjektit/views/components"
	"musiikkiProjektit/views/notes"
	"net/http"
)

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

func HandleCreateNewNote(db *sql.DB, w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie(SESSION_TOKEN_NAME)
	if err != nil {
		//TODO: joku parempi virheellinen session token vastaus
		w.WriteHeader(401)
		return
	}

	userSession := Sessions[cookie.Value]
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
	newNote := utils.Note{
		Title: r.FormValue("newNoteTitle"),
		Note: r.FormValue("newNote"),
	}

	err, errors := parseNewNote(db, newNote, userSession.Username)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	err = components.NewNoteForm(newNote, errors).Render(r.Context(), w)
	if err != nil {
		log.Println(err)
	}
}

func parseNewNote(db *sql.DB, noteData utils.Note, username string) (error,[]string){
	insertQuery := `
	INSERT INTO notes(user_id, title, note, created)
	VALUES((SELECT id FROM users where username = ?),?,?, datetime('now'));
	`
	var errors []string
	errors = append(errors, validateNewNoteTitle(noteData.Title)...)
	errors = append(errors, validateNewNoteText(noteData.Note)...)
	if len(errors) > 0{
		return nil, errors
	}

	_, err := db.Exec(insertQuery, username, noteData.Title, noteData.Note)
	if err != nil {
		return err, nil
	}
	return nil, errors
}

func validateNewNoteTitle(title string) []string{
	var errors []string
	//TODO: joku parempi validointi
	if len(title) < 3{
		errors = append(errors, "Title too short. It should be atleast 3 characters")
	}
	return errors
}
func validateNewNoteText(text string) []string{
	var errors []string
	//TODO: joku parempi validointi
	if len(text) < 5{
		errors = append(errors, "Note too short. It should be atleast 5 characters")
	}
	return errors
}
