package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"musiikkiProjektit/auth"
	"musiikkiProjektit/handlers"
	"musiikkiProjektit/session"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// dbConnect establishes the database connection
//
// returns a pointer to the database and nil,
// or nil and error if the connection could not be established
func dbConnect() (*sql.DB, error){
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// initializeDBSchema initializes the database schema on startup
// if it does not exist
//
// requires a database pointer, and return error if initialization was not successful
func initializeDBSchema(db *sql.DB) error{
	err := initializeUsersSchema(db)
	if err != nil {
		return fmt.Errorf("%s, while creating the users schema", err)
	}

	err = initializeNotesSchema(db)
	if err != nil {
		return fmt.Errorf("%s, while creating the notes schema", err)
	}

	err = initializeSessionStorage(db)
	if err != nil {
		return fmt.Errorf("%s, while creating the session storage", err)
	}

	return nil
}

// initializeUsersSchema creates the schema for saving users,
// if the table does not exist
//
// requires a database pointer, and return error if initialization was not successful
func initializeUsersSchema(db *sql.DB) error{
	initQuery := `CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`

	_, err := db.Exec(initQuery)
	if err != nil {
		return err
	}

	return nil
}

// initializeNotesSchema creates the schema for saving notes,
// if the table does not exist
//
// requires a database pointer, and return error if initialization was not successful
func initializeNotesSchema(db *sql.DB) error{
	initQuery := `CREATE TABLE IF NOT EXISTS notes(
		note_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		title TEXT NOT NULL,
		note BLOB NOT NULL,
		created TEXT NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err := db.Exec(initQuery)
	if err != nil {
		return err
	}

	return nil
}

// initializeSessionStorage creates the schema for saving session data,
// if the table does not exist
//
// requires a database pointer, and return error if initialization was not successful
func initializeSessionStorage(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS sessions(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		ip TEXT NOT NULL,
		timestamp TEXT NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// configs is a struct that contains the configurations that can be done outside
// the program
var configs struct{
	Port int `json:"port"`
}

// readConfigs reads the configurations from the config.json file in the current
// directory, to a configs struct.
//
// returns error if reading was not successful
func readConfigs() error{
	file, err := os.Open("./config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configs)
	if err != nil{
		return err
	}

	return nil
}

func init(){
	err := readConfigs()
	if err != nil {
		log.Fatalf("Failed to read the configuration file %s\n", err)
	}

	go session.CleanupOutdatedSessions(30 * time.Minute)
}

func main(){
	handler := http.NewServeMux()
	server := http.Server{
		Addr: ":" + strconv.Itoa(configs.Port),
		Handler: handler,
		ReadTimeout: 30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout: 120 * time.Second,
	}

	db, err := dbConnect()
	if err != nil {
		log.Fatalf("Error establishing the database connectin: %s\n", err)
	}
	defer db.Close()

	err = initializeDBSchema(db)
	if err != nil {
		log.Fatalf("Error creating the database schema: %s\n", err)
	}

	// for testing ====
	testUserName := "test"
	testUserPassword := "123"
	err = auth.CreateNewUser(db, testUserName, testUserPassword )
	if err != nil {
		if err.Error() != "UNIQUE constraint failed: users.username"{
			log.Fatalln(err)
		}
	}
	//===

	// Pages
	handler.HandleFunc("GET /", session.HandleSessionMiddleware(handlers.HandleServeIndex))
	handler.HandleFunc("GET /notes", session.HandleSessionMiddleware(handlers.HandleServeNotes))
	handler.HandleFunc("GET /chordprogress", session.HandleSessionMiddleware(handlers.HandleServeChordProg))
	handler.HandleFunc("GET /keyquiz", session.HandleSessionMiddleware(handlers.HandleServeKeyQuiz))
	handler.HandleFunc("GET /login", session.HandleSessionMiddleware(handlers.HandleLoginPage))
	handler.HandleFunc("GET /oldchordprog", session.HandleSessionMiddleware(handlers.HandleServeOldChordProg))



	/*___API___*/
	// keyquiz
	handler.HandleFunc("GET /api/keyquiz/start", handlers.HandleStartKeyQuiz)
	handler.HandleFunc("POST /api/keyquiz/checkanswer", handlers.HandleCheckQuiz)
	handler.HandleFunc("GET /api/keyquiz/newquiz", handlers.HandleStartKeyQuiz)


	// Session
	handler.HandleFunc("POST /api/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleLogin(db, w,r)
	})
	handler.HandleFunc("GET /logout", handlers.HandleLogout)


	// Notes
	handler.HandleFunc("POST /api/notes", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleCreateNewNote(db, w, r)
	})
	handler.HandleFunc("GET /notes/form", handlers.HandleNewNoteForm)
	handler.HandleFunc("GET /api/notes", func(w http.ResponseWriter, r *http.Request) { 
		handlers.HandleGetSavedNotes(db, w, r)
	})
	handler.HandleFunc("DELETE /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleDeleteNote(db, w, r)
	})
	handler.HandleFunc("GET /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleEditNote(db, w ,r)
	})


	// Files
	handler.Handle("GET /index.js", http.FileServer(http.Dir("./public/")))
	handler.Handle("GET /old-chord-prog.js", http.FileServer(http.Dir("./public/")))


	log.Printf("server started on port %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}
