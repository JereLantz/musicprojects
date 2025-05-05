package handlers

import (
	"musiikkiProjektit/views/login"
	"net/http"
)

func HandleLoginPage(w http.ResponseWriter, r *http.Request){
	login.LoginPage().Render(r.Context(), w)
}
