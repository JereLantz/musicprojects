package handlers

import (
	"database/sql"
	"log"
	"musiikkiProjektit/views/login"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func HandleLoginPage(w http.ResponseWriter, r *http.Request){
	login.LoginPage().Render(r.Context(), w)
}


func HandleLogin(db *sql.DB, w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(500)
		log.Printf("failed login because could not parse form: %s\n", err)
		return
	}

	username := r.FormValue("login-uname")
	password := r.FormValue("login-passwd")
	err = checkUserCredentials(db, username, password)

	if err != nil {
		log.Println("failed")
		w.WriteHeader(400)
	}else{
		log.Println("Success")
	}
}

func checkUserCredentials(db *sql.DB, uname, passwd string) error{
	fetchPasswdQuery := `SELECT password FROM users WHERE username = ?;`
	var storedPasswd []byte

	passRow := db.QueryRow(fetchPasswdQuery, uname)

	err := passRow.Scan(&storedPasswd)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(storedPasswd, []byte(passwd))

	return err
}
