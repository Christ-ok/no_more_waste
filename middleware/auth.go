package middleware

import (
	"net/http"
	"no_more_waste/session"
)

const nomSession = "nmw-session"

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		sess, sessErr := session.Store.Get(request, nomSession)
		if sessErr != nil {
			http.Error(response, "Erreur lors de la récupération du cookie", http.StatusNotFound)
			return
		}

		authentifie, ok := sess.Values["authentifie"].(bool)
		if !ok || !authentifie {
			http.Redirect(response, request, "/connexion", http.StatusFound)
			return
		}

		next(response, request)
	}
}

func AuthRole(next http.HandlerFunc, rolesAutorises ...string) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		sess, sessErr := session.Store.Get(request, nomSession)
		if sessErr != nil {
			http.Error(response, "Erreur lors de la récupération du cookie", http.StatusNotFound)
			return
		}

		authentifie, ok := sess.Values["authentifie"].(bool)
		if !ok || !authentifie {
			http.Redirect(response, request, "/connexion", http.StatusFound)
			return
		}

		roleUtilisateur, ok := sess.Values["role"].(string)
		if !ok {
			http.Redirect(response, request, "/connexion", http.StatusFound)
			return
		}

		if !roleIsAllowed(roleUtilisateur, rolesAutorises) {
			http.Error(response, "Accès refusé", http.StatusForbidden)
			return
		}

		next(response, request)
	}
}

func roleIsAllowed(role string, rolesAutorises []string) bool {
	for _, roleAllowed := range rolesAutorises {
		if roleAllowed == role {
			return true
		}
	}

	return false
}

func Logout(response http.ResponseWriter, request *http.Request) {
	sess, _ := session.Store.Get(request, nomSession)
	sess.Options.MaxAge = -1
	sess.Save(request, response)
	http.Redirect(response, request, "/", http.StatusFound)
}
