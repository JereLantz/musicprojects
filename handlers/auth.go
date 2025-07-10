package handlers

import (
	"database/sql"
	"log"
	"musiikkiProjektit/session"
	"musiikkiProjektit/views/components"
	"musiikkiProjektit/views/login"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

func HandleLoginPage(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie(session.SessionTokenName)
	if err != nil{
		w.WriteHeader(200)
		login.LoginPage(session.Session{}).Render(r.Context(), w)
		return
	}
	sessionData, err := session.GetSession(cookie.Value)
	if err != nil {
		w.WriteHeader(200)
		login.LoginPage(session.Session{}).Render(r.Context(), w)
		return
	}
	w.WriteHeader(200)
	login.LoginPage(sessionData).Render(r.Context(), w)
}

func HandleLogin(db *sql.DB, w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(500)
		log.Printf("failed login because could not parse form: %s\n", err)
		return
	}

	inputtedUsername := r.FormValue("login-uname")
	inputtedPassword := r.FormValue("login-passwd")
	err = checkUserCredentials(db, inputtedUsername, inputtedPassword)

	if err == nil {
		log.Printf("Successfull login %s\n", inputtedUsername)
		err = session.SessionLogin(r, inputtedUsername)
		if err != nil {
			w.WriteHeader(500)
			log.Printf("error modifying the session with login details: %s\n", err)
			return
		}
		w.Header().Add("Hx-Redirect", "/")
		return
	}

	log.Printf("Failed login, %s\n", err)
	components.LoginForm(true, inputtedUsername, inputtedPassword).Render(r.Context(), w)
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

func HandleLogout(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie(session.SessionTokenName)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Failed get the session cookie to logout %s\n", err)
		return
	}

	session.DeleteSession(cookie.Value)
	
	http.SetCookie(w, &http.Cookie{
		Name: session.SessionTokenName,
		Value: "",
		Expires: time.Now(),
	})
	http.Redirect(w,r,"/", 303)
}

func CreateNewUser(db *sql.DB, username, password string) error{
	createNewUserQuery := `INSERT INTO users(username, password) VALUES(?,?);`
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(createNewUserQuery, username, passwordHash)
	if err != nil {
		return err
	}
	return nil
}
