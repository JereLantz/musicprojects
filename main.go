package main

import (
	"database/sql"
	"fmt"
	"log"
	"musiikkiProjektit/handlers"
	"musiikkiProjektit/utils"
	"musiikkiProjektit/views/index"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func handleServeIndex(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		w.WriteHeader(404)
		return
	}
	cookie, err := r.Cookie(handlers.SESSION_TOKEN_NAME)
	if err != nil{
		log.Printf("Failed to fetch the session for displaying the home. %s\n", err)
		w.WriteHeader(500)
		return
	}
	sessionData := handlers.Sessions[cookie.Value]

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

func init(){
	go handlers.CleanupOutdatedSessions(30 * time.Minute)
}

func main(){
	handler := http.NewServeMux()
	server := http.Server{
		Addr: ":42069",
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
	handler.HandleFunc("GET /", handlers.HandleSessionMiddleware(handleServeIndex))
	handler.HandleFunc("GET /notes", handlers.HandleSessionMiddleware(handlers.HandleServeNotes))
	handler.HandleFunc("GET /chordprogress", handlers.HandleSessionMiddleware(handlers.HandleServeChordProg))
	handler.HandleFunc("GET /keyquiz", handlers.HandleSessionMiddleware(handlers.HandleServeKeyQuiz))
	handler.HandleFunc("GET /login", handlers.HandleSessionMiddleware(handlers.HandleLoginPage))

	// API
	handler.HandleFunc("GET /api/keyquiz/start", handlers.HandleStartKeyQuiz)
	handler.HandleFunc("POST /api/keyquiz/checkanswer", handlers.HandleCheckQuiz)
	handler.HandleFunc("GET /api/keyquiz/newquiz", handlers.HandleStartKeyQuiz)
	handler.HandleFunc("POST /api/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleLogin(db, w,r)
	})
	handler.HandleFunc("GET /logout", handlers.HandleLogout)
	handler.HandleFunc("POST /api/notes/addnew", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleCreateNewNote(db, w, r)
	})

	// Files
	handler.Handle("GET /index.js", http.FileServer(http.Dir("./")))

	log.Printf("server started on port %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}
