package notes

import "database/sql"

// Note represents the structure of a saved note
type Note struct {
	Id int
	Title string
	Note string
	//TODO: lisää timestamp?
}

// GetUsersNotes returns a slice of Note structs and an error
//
// Requires db pointer and the users username
func GetUsersNotes(db *sql.DB, username string) ([]Note, error){
	var userNotes []Note
	//TODO: hae myös time stamp?
	query := `SELECT note_id, title, note FROM NOTES WHERE user_id = (SELECT id FROM users WHERE username = ?);`

	row, err := db.Query(query, username)
	if err != nil {
		return []Note{}, err
	}
	defer row.Close()

	for row.Next(){
		var note Note
		err = row.Scan(&note.Id, &note.Title, &note.Note)
		if err != nil {
			return []Note{}, err
		}
		userNotes = append(userNotes, note)
	}

	return userNotes, nil
}

// TODO: muuta näitä siten että on 1 validate field functio ja tallenna funktio
func ParseNewNote(db *sql.DB, noteData Note, username string) (error,[]string){
	insertQuery := `
	INSERT INTO notes(user_id, title, note, created)
	VALUES((SELECT id FROM users where username = ?),?,?, datetime('now'));
	`
	var errors []string
	errors = append(errors, ValidateNewNoteTitle(noteData.Title)...)
	errors = append(errors, ValidateNewNoteText(noteData.Note)...)
	if len(errors) > 0{
		return nil, errors
	}

	_, err := db.Exec(insertQuery, username, noteData.Title, noteData.Note)
	if err != nil {
		return err, nil
	}
	return nil, errors
}

func ValidateNewNoteTitle(title string) []string{
	var errors []string
	//TODO: joku parempi validointi
	if len(title) < 3{
		errors = append(errors, "Title too short. It should be atleast 3 characters")
	}
	return errors
}

func ValidateNewNoteText(text string) []string{
	var errors []string
	//TODO: joku parempi validointi
	if len(text) < 5{
		errors = append(errors, "Note too short. It should be atleast 5 characters")
	}
	return errors
}
