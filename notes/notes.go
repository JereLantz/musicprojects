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
// Return the created id and an error if not successful
func (n Note)SaveNewNote(db *sql.DB, username string) (int, error){
	now := time.Now()
	insertQuery := `
	INSERT INTO notes(user_id, title, note, created)
	VALUES((SELECT id FROM users where username = ?),?,?, ?);
	`
	result, err := db.Exec(insertQuery, username, n.Title, n.Note, now.Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
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


// DeleteNote deletes a note with the specified id from the user with the
// specified username
//
// Returns error if unsuccessful
func DeleteNote(db *sql.DB, noteID int, username string) error {
	query := `DELETE FROM notes WHERE note_id = ? AND user_id = (SELECT id FROM users WHERE username = ?);`

	_, err := db.Exec(query, noteID, username)
	if err != nil {
		return err
	}
	return nil
}


// GetNote Fetches a specific note from the database using the notes id an username
//
// returns Note struct with the data and nil, or empty Note struct and an error
func GetNote(db *sql.DB, noteID int, username string) (Note, error) {
	var note Note
	var timeStamp string
	query := `SELECT note_id, title, note, created FROM notes WHERE note_id =? AND user_id = (SELECT id FROM users WHERE username = ?);`

	row := db.QueryRow(query, noteID, username)

	err := row.Scan(&note.Id, &note.Title, &note.Note, &timeStamp)
	if err != nil {
		return Note{}, err
	}

	note.Created, err = time.Parse(time.RFC3339, timeStamp)
	if err != nil {
		return Note{}, err
	}

	return note, nil
}
