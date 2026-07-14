package session

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

func Init() {
	Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Domain:   "",
	}
}
