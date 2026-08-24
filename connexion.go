package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"no_more_waste/i18n"
	"no_more_waste/middleware"
	"no_more_waste/models"
	"no_more_waste/session"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func Login(database *sql.DB) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		var user models.Utilisateur

		if errForm := request.ParseForm(); errForm != nil {
			fmt.Printf("Erreur Parse Form : %v", errForm)
		}

		user.Email = request.FormValue("email")
		motDePasseSaisi := request.FormValue("mot_de_passe")

		if user.Email == "" || motDePasseSaisi == "" {
			http.Error(response, "Email et mot de passe non saisi", http.StatusBadRequest)
			return
		}

		var roleUser string

		errRole := database.QueryRow(`SELECT u.id_utilisateur, u.nom, u.prenom, u.mot_de_passe,
			        					u.telephone, u.adresse, u.ville, u.code_postal, u.pays,
			        					r.nom
			 							FROM utilisateur u
			 							JOIN role r ON r.id_role = u.id_role
			 							WHERE u.email = $1`, user.Email).Scan(&user.IDUtilisateur, &user.Nom, &user.Prenom, &user.MotDePasse,
			&user.Telephone, &user.Adresse, &user.Ville, &user.CodePostal, &user.Pays,
			&roleUser,
		)

		if errors.Is(errRole, sql.ErrNoRows) {
			http.Error(response, "Email ou mot de passe incorrect", http.StatusUnauthorized)
			return
		}

		if errRole != nil {
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}

		if bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.MotDePasse), []byte(motDePasseSaisi)); bcryptErr != nil {
			http.Error(response, "Email ou mot de passe incorrect", http.StatusUnauthorized)
			return
		}

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			fmt.Printf("ERREUR SESSION GET : %v\n", sessErr)
			http.Error(response, "Erreur lors de la création de la session", http.StatusInternalServerError)
			return
		}

		sess.Values["authentifie"] = true
		sess.Values["id_utilisateur"] = user.IDUtilisateur
		sess.Values["role"] = roleUser
		sess.Values["nom"] = user.Nom
		sess.Values["prenom"] = user.Prenom

		if err := sess.Save(request, response); err != nil {
			fmt.Printf("Erreur lors de la sauvegarde la session : %v", err)
			http.Error(response, "Erreur lors de la sauvegarde de la session", http.StatusInternalServerError)
			return
		}

		http.Redirect(response, request, RedirectionSelonRole(roleUser), http.StatusFound)
	}
}

func RedirectionSelonRole(roleNom string) string {
	switch roleNom {
	case "ADMIN_GENERAL":
		return "/admin"

	case "ADMIN_AGENCE":
		return "/admin-agence"

	case "COMMERCANT":
		return "/commercant"

	case "BENEVOLE":
		return "/benevole"

	case "ADHERENT":
		return "/adherent"

	case "ASSOCIATION":
		return "/association"

	default:
		return "/"
	}
}

func PageConnexion(response http.ResponseWriter, request *http.Request) {
	language := middleware.GetLanguage(request)
	fmt.Println("LANGUE ACTUELLE :", language)
	fmt.Println("TRADUCTION TEST :", i18n.Traduction(language, "login.title"))

	tmpl, err := template.New("connexion.html").Funcs(template.FuncMap{
		"t": func(key string) string {
			return i18n.Traduction(language, key)
		},
	}).
		ParseFiles("templates/connexion.html")

	if err != nil {
		http.Error(response, "Erreur lors du chargement de la page", http.StatusInternalServerError)
		fmt.Println("Erreur template connexion :", err)
		return
	}

	err = tmpl.ExecuteTemplate(response, "connexion.html", nil)

	if err != nil {
		fmt.Println("Erreur exécution template connexion :", err)
	}
}

func PageBenevole(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "templates/benevole/accueil_Benevole.html")
}

func PageAdminGeneral(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "templates/admin_general/adminGeneral_accueil.html")
}

func PageAdminAgence(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "templates/admin_agence/adminAgence_accueil.html")
}

func PageCommercant(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "templates/commercant/commercant_accueil.html")
}

func PageAdherent(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "templates/adherent/adherent_accueil.html")
}

func PageAssociation(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "templates/association/association_accueil.html")
}
