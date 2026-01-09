package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Bronku/cake-order-tracker/models"
)

func (a *Authenticator) getSession(r *http.Request) (models.Token, error) {
	c, err := r.Cookie("token")
	if err != nil {
		return models.Token{}, err
	}
	session, ok := a.sessions[c.Value]
	if !ok {
		return models.Token{}, errors.New("session not found")
	}
	if time.Since(session.Expiration) > 0 {
		delete(a.sessions, c.Value)
		err := a.s.CleanSessions()
		if err != nil {
			fmt.Println("error cleaning the sessions", err)
		}
		return models.Token{}, errors.New("session expired")
	}
	return session, nil
}

func (a *Authenticator) login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	login := r.PostFormValue("login")
	password := r.PostFormValue("password")
	if a.verifyCredentials(login, password) != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	cookie, err := a.newSession(login)
	if err != nil {
		w.Header().Set("content-type", "text/html")
		fmt.Fprint(w, "internal server error")
		return
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *Authenticator) logout(w http.ResponseWriter, r *http.Request) {
	var cookie http.Cookie
	cookie.Name = "token"
	cookie.Value = "nil"
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Path = "/"
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
	c, err := r.Cookie("token")
	if err != nil {
		return
	}
	_, ok := a.sessions[c.Value]
	if !ok {
		return
	}
	fmt.Println("removing session")
	delete(a.sessions, c.Value)
	err = a.s.RevokeSession(c.Value)
	fmt.Println(err)
}
