package notes

import (
	"database/sql"
	"time"
)

// Note represents the structure of a saved note
type Note struct {
	Id int
	Title string
	Note string
	Created time.Time
}

// GetUsersNotes returns a slice of Note structs and an error
//
// The notes are ordered in a descending order, based on note id
//
// Requires db pointer and the users username
func GetUsersNotes(db *sql.DB, username string) ([]Note, error){
	var userNotes []Note
	query := `SELECT note_id, title, note, created
	FROM NOTES WHERE user_id = (SELECT id FROM users WHERE username = ?)
	ORDER BY note_id DESC;`

	row, err := db.Query(query, username)
	if err != nil {
		return []Note{}, err
	}
	defer row.Close()

	for row.Next(){
		var note Note
		var timeStr string
		err = row.Scan(&note.Id, &note.Title, &note.Note, &timeStr)
		if err != nil {
			return []Note{}, err
		}
		note.Created, err = time.Parse(time.RFC3339, timeStr)
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
	now := time.Now()
	insertQuery := `
	INSERT INTO notes(user_id, title, note, created)
	VALUES((SELECT id FROM users where username = ?),?,?, ?);
	`
	_, err := db.Exec(insertQuery, username, n.Title, n.Note, now.Format(time.RFC3339))
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
