package handlers

import (
	"database/sql"
	"log"
	"musiikkiProjektit/utils"
	"musiikkiProjektit/views/components"
	"musiikkiProjektit/views/login"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const BCRYPT_COST = 12

func HandleLoginPage(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie(SESSION_TOKEN_NAME)
	if err != nil{
		log.Printf("Failed to fetch the session for displaying the login page. %s\n", err)
		w.WriteHeader(500)
		return
	}
	sessionData := Sessions[cookie.Value]
	login.LoginPage(sessionData).Render(r.Context(), w)
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

	if err == nil {
		log.Printf("Successfull login %s\n", inputtedCreds.Username)
		err = SessionLogin(r, inputtedCreds)
		if err != nil {
			w.WriteHeader(500)
			log.Printf("error modifying the session with login details: %s\n", err)
			return
		}
		w.Header().Add("Hx-Redirect", "/")
		return
	}

	log.Printf("Failed login, %s\n", err)
	components.LoginForm(true, inputtedCreds).Render(r.Context(), w)
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

func HandleLogout(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie(SESSION_TOKEN_NAME)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Failed get the session cookie to logout %s\n", err)
		return
	}

	delete(Sessions, cookie.Value)
	
	http.SetCookie(w, &http.Cookie{
		Name: SESSION_TOKEN_NAME,
		Value: "",
		Expires: time.Now(),
	})
	http.Redirect(w,r,"/", 303)
}

func CreateNewUser(db *sql.DB, credentials utils.Credentials) error{
	createNewUserQuery := `INSERT INTO users(username, password) VALUES(?,?);`
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), BCRYPT_COST)
	if err != nil {
		return err
	}

	_, err = db.Exec(createNewUserQuery, credentials.Username, passwordHash)
	if err != nil {
		return err
	}
	return nil
}
