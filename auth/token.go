package auth

import (
	"net/http"
	"time"

	"github.com/Bronku/cake-order-tracker/crypto"
	"github.com/Bronku/cake-order-tracker/models"
)

func (a *Authenticator) newSession(user string) (http.Cookie, error) {
	key := crypto.GenerateKey()
	var session models.Token
	var cookie http.Cookie
	session.User = user
	session.Expiration = time.Now().Add(time.Hour * 24)
	a.sessions[key] = session
	err := a.s.AddSession(key, user, session.Expiration)
	if err != nil {
		return cookie, err
	}

	cookie.Name = "token"
	cookie.Value = key
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Path = "/"
	//cookie.Secure = true
	return cookie, nil
}
