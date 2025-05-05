package main

import (
	"database/sql"
	"fmt"
	"log"
	"musiikkiProjektit/handlers"
	"musiikkiProjektit/views/index"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func handleServeIndex(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(200)
	index.Index().Render(r.Context(), w)
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
		username text NOT NULL,
		salt text NOT NULL,
		password text NOT NULL
	);`

	_, err := db.Exec(initQuery)
	if err != nil {
		return err
	}

	return nil
}

func initializeNotesSchema(db *sql.DB) error{
	initQuery := `CREATE TABLE IF NOT EXISTS notes(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title text NOT NULL,
		note text NOT NULL
	);`

	_, err := db.Exec(initQuery)
	if err != nil {
		return err
	}

	return nil
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

	// Pages
	handler.HandleFunc("GET /", handleServeIndex)
	handler.HandleFunc("GET /notes", handlers.HandleServeNotes)
	handler.HandleFunc("GET /chordprogress", handlers.HandleServeChordProg)
	handler.HandleFunc("GET /keyquiz", handlers.HandleServeKeyQuiz)
	handler.HandleFunc("GET /login", handlers.HandleLoginPage)

	// API
	handler.HandleFunc("POST /keyquiz/start", handlers.HandleStartKeyQuiz)
	handler.HandleFunc("POST /keyquiz/checkanswer", handlers.HandleCheckQuiz)
	handler.HandleFunc("POST /keyquiz/newquiz", handlers.HandleStartKeyQuiz)

	// Files
	handler.Handle("GET /index.js", http.FileServer(http.Dir("./")))

	log.Printf("server started on port %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}
