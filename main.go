package main

import (
	"log"
	"net/http"
)

func main(){
	handler := http.NewServeMux()
	server := http.Server{
		Addr: ":42069",
		Handler: handler,
	}

	log.Printf("server started on port %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}
