package auth

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"

	"github.com/Bronku/cake-order-tracker/models"
	"github.com/Bronku/cake-order-tracker/store"
)

type Authenticator struct {
	sessions map[string]models.Token
	s        *store.Store
}

func New(s *store.Store) *Authenticator {
	var out Authenticator
	var err error
	out.s = s
	out.sessions, err = s.GetSessions()
	//fmt.Println(out.sessions)
	if err != nil {
		log.Fatal(err)
	}
	return &out
}

func (a *Authenticator) ensureAuth(in http.Handler) http.Handler {
	fmt.Println("ensureAuth called")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := a.getSession(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		in.ServeHTTP(w, r)
	})
}

func (a *Authenticator) Middleware(in http.Handler) http.Handler {
	handler := http.NewServeMux()
	handler.HandleFunc("GET /login", getLogin)
	handler.HandleFunc("POST /login", a.login)
	handler.HandleFunc("GET /logout", a.logout)
	handler.Handle("/", a.ensureAuth(in))
	return handler
}
