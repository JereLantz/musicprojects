package main

import (
	"log"
	"musiikkiProjektit/views/chordProgress"
	"musiikkiProjektit/views/index"
	"musiikkiProjektit/views/keyQuiz"
	"musiikkiProjektit/views/notes"
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

func handleServeNotes(w http.ResponseWriter, r *http.Request){
	notes.NotesPage().Render(r.Context(), w)
}

func handleServeChordProg(w http.ResponseWriter, r *http.Request){
	chordprogress.ChordProgPage().Render(r.Context(), w)
}

func handleServeKeyQuiz(w http.ResponseWriter, r *http.Request){
	keyquiz.KeyQuizPage().Render(r.Context(), w)
}

func main(){
	handler := http.NewServeMux()
	server := http.Server{
		Addr: ":42069",
		Handler: handler,
	}

	// Pages
	handler.HandleFunc("GET /", handleServeIndex)
	handler.HandleFunc("GET /notes", handleServeNotes)
	handler.HandleFunc("GET /chordprogress", handleServeChordProg)
	handler.HandleFunc("GET /keyquiz", handleServeKeyQuiz)

	// API

	// Files

	log.Printf("server started on port %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}
