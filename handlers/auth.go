package handlers

import (
	"database/sql"
	"log"
	"musiikkiProjektit/utils"
	"musiikkiProjektit/views/components"
	"musiikkiProjektit/views/login"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func HandleLoginPage(w http.ResponseWriter, r *http.Request){
	login.LoginPage().Render(r.Context(), w)
}


func HandleLogin(db *sql.DB, w http.ResponseWriter, r *http.Request){
	var inputtedCreds utils.Credentials
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(500)
		log.Printf("failed login because could not parse form: %s\n", err)
		return
	}

	inputtedCreds.Username = r.FormValue("login-uname")
	inputtedCreds.Password = r.FormValue("login-passwd")
	err = checkUserCredentials(db, inputtedCreds)

	if err != nil {
		log.Printf("Failed login, %s\n", err)
		components.LoginForm(true, inputtedCreds).Render(r.Context(), w)
	}else{
		log.Printf("Successfull login %s\n", inputtedCreds.Username)
		w.Header().Add("Hx-Retarget", "#main-content")
		components.LoginWelcomeMsg(inputtedCreds.Username).Render(r.Context(), w)
	}
}

func checkUserCredentials(db *sql.DB, credentials utils.Credentials) error{
	fetchPasswdQuery := `SELECT password FROM users WHERE username = ?;`
	var storedPasswd []byte

	passRow := db.QueryRow(fetchPasswdQuery, credentials.Username)

	err := passRow.Scan(&storedPasswd)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(storedPasswd, []byte(credentials.Password))

	return err
}
