package main

import (
	"log"
	"musiikkiProjektit/views/index"
	"net/http"
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

	handler.HandleFunc("GET /", handleServeIndex)
	log.Printf("server started on port %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}
