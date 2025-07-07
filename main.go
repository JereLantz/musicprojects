package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"musiikkiProjektit/handlers"
	"musiikkiProjektit/session"
	"musiikkiProjektit/utils"
	"musiikkiProjektit/views/index"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func handleServeIndex(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		w.WriteHeader(404)
		return
	}
	cookie, err := r.Cookie(session.SessionTokenName)
	if err != nil{
		w.WriteHeader(200)
		index.Index(session.Session{}).Render(r.Context(), w)
		return
	}
	sessionData, err := session.GetSession(cookie.Value)
	if err != nil {
		index.Index(session.Session{}).Render(r.Context(), w)
		log.Printf("could not get session information when serving the index %s\n", err)
		return
	}

	w.WriteHeader(200)
	index.Index(sessionData).Render(r.Context(), w)
}

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

func initializeSessionStorage(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS sessions(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT UNIQUE NOT NULL,
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

type config struct{
	Port int `json:"port"`
}

var configs config

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
	}

	db, err := dbConnect()
	defer db.Close()
	if err != nil {
		log.Fatalf("Error establishing the database connectin: %s\n", err)
	}

	err = initializeDBSchema(db)
	if err != nil {
		log.Fatalf("Error creating the database schema: %s\n", err)
	}

	// for testing ====
	testUser := utils.Credentials{
		Username: "test",
		Password: "123",
	}
	err = handlers.CreateNewUser(db, testUser)
	if err != nil {
		if err.Error() != "UNIQUE constraint failed: users.username"{
			log.Fatalln(err)
		}
	}
	//===

	// Pages
	handler.HandleFunc("GET /", session.HandleSessionMiddleware(handleServeIndex))
	handler.HandleFunc("GET /notes", session.HandleSessionMiddleware(handlers.HandleServeNotes))
	handler.HandleFunc("GET /chordprogress", session.HandleSessionMiddleware(handlers.HandleServeChordProg))
	handler.HandleFunc("GET /keyquiz", session.HandleSessionMiddleware(handlers.HandleServeKeyQuiz))
	handler.HandleFunc("GET /login", session.HandleSessionMiddleware(handlers.HandleLoginPage))



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


	// Files
	handler.Handle("GET /index.js", http.FileServer(http.Dir("./public/")))


	log.Printf("server started on port %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}
