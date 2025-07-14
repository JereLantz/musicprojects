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
// The notes are ordered in a descending order, based on note id
//
// Requires db pointer and the users username
func GetUsersNotes(db *sql.DB, username string) ([]Note, error){
	var userNotes []Note
	//TODO: hae myös time stamp?
	query := `SELECT note_id, title, note
	FROM NOTES WHERE user_id = (SELECT id FROM users WHERE username = ?)
	ORDER BY note_id DESC;`

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


// SaveNewNote creates new note from the note struct for the user with the
// supplied username in supplied database.
//
// Return error if not successful
func (n Note)SaveNewNote(db *sql.DB, username string) (error){
	insertQuery := `
	INSERT INTO notes(user_id, title, note, created)
	VALUES((SELECT id FROM users where username = ?),?,?, datetime('now'));
	`
	_, err := db.Exec(insertQuery, username, n.Title, n.Note)
	if err != nil {
		return err
	}
	return nil
}


// ValidateNoteFields checks that the fields contain valid data
//
// returns a string array that contains the possible user input errors,
// and error for possible parsing errors
func (n Note)Validate() ([]string, error){
	var userErrors []string
	//TODO:
	return userErrors, nil
}
