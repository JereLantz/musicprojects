package handlers

import (
	"database/sql"
	"log"
	"musiikkiProjektit/auth"
	"musiikkiProjektit/session"
	"musiikkiProjektit/views/components"
	"musiikkiProjektit/views/pages"
	"net/http"
	"net/url"
	"time"
)

func HandleLoginPage(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie(session.SessionTokenName)
	if err != nil{
		w.WriteHeader(200)
		pages.Login(session.Session{}).Render(r.Context(), w)
		return
	}
	sessionData, err := session.GetSession(cookie.Value)
	if err != nil {
		w.WriteHeader(200)
		pages.Login(session.Session{}).Render(r.Context(), w)
		return
	}
	w.WriteHeader(200)
	pages.Login(sessionData).Render(r.Context(), w)
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
	err = auth.CheckUserCredentials(db, inputtedUsername, inputtedPassword)

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

// HandleLogout handles logging the user out. Then redirects the user
// to the home page.
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
	requestUrl := r.Header.Get("referer")
	parsedUrl, err := url.Parse(requestUrl)
	if err != nil {
		http.Redirect(w,r, "/", 303)
		return
	}

	switch parsedUrl.Path {
	case "/keyquiz":
		fallthrough
	case "/chordprogress":
		fallthrough
	case "/notes":
		http.Redirect(w,r, requestUrl, 303)
	default:
		http.Redirect(w,r, "/", 303)
	}
}
