package main

import (
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

func main(){
	handler := http.NewServeMux()
	server := http.Server{
		Addr: ":42069",
		Handler: handler,
	}

	// Pages
	handler.HandleFunc("GET /", handleServeIndex)
	handler.HandleFunc("GET /notes", handlers.HandleServeNotes)
	handler.HandleFunc("GET /chordprogress", handlers.HandleServeChordProg)
	handler.HandleFunc("GET /keyquiz", handlers.HandleServeKeyQuiz)

	// API
	handler.HandleFunc("POST /keyquiz/start", handlers.HandleStartKeyQuiz)

	// Files

	log.Printf("server started on port %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}
