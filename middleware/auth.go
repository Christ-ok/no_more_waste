package middleware

import (
	"database/sql"
	"fmt"
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

func GetIDAgenceUtilisateur(database *sql.DB, idUtilisateur int) (int, error) {
	var idAgence sql.NullInt64

	err := database.QueryRow(
		"SELECT id_agence FROM utilisateur WHERE id_utilisateur = $1",
		idUtilisateur,
	).Scan(&idAgence)
	if err != nil {
		return 0, err
	}

	if !idAgence.Valid {
		return 0, fmt.Errorf("aucune agence associée à l'utilisateur %d", idUtilisateur)
	}
	return int(idAgence.Int64), nil
}

func Logout(response http.ResponseWriter, request *http.Request) {
	sess, _ := session.Store.Get(request, nomSession)
	sess.Options.MaxAge = -1
	if err := sess.Save(request, response); err != nil {
		fmt.Printf("Erreur déconnexion : %v", err)
	}
	http.Redirect(response, request, "/", http.StatusFound)
}
