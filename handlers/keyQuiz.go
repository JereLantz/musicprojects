package handlers

import (
	"musiikkiProjektit/views/components"
	"musiikkiProjektit/views/keyQuiz"
	"net/http"
)

func HandleServeKeyQuiz(w http.ResponseWriter, r *http.Request){
	keyquiz.KeyQuizPage().Render(r.Context(), w)
}

func HandleStartKeyQuiz(w http.ResponseWriter, r *http.Request){
	components.KeyQuiz().Render(r.Context(), w)
}
