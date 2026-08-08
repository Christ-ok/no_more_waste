package routes

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"no_more_waste/middleware"
	"no_more_waste/models"
	"no_more_waste/session"

	"golang.org/x/crypto/bcrypt"
)

type CommercantData struct {
	Commercant     models.Utilisateur
	Nom_Entreprise string
	Type_Commerce  string
	Numero_Siret   string
}

func AfficherPageModifierProfilCommercant(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur récupération session", http.StatusInternalServerError)
			return
		}

		idUtilisateur, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			http.Error(response, "Erreur récupération ID", http.StatusInternalServerError)
			return
		}

		var user models.Utilisateur
		var nom_entreprise, type_commerce, numero_siret string

		rowsErr := database.QueryRow(`SELECT u.nom, u.prenom, u.email, u.telephone, u.adresse, u.ville, u.code_postal, u.pays, c.nom_entreprise, c.type_commerce, c.numero_siret FROM utilisateur u
										JOIN commercant c ON c.id_utilisateur = u.id_utilisateur
									  WHERE u.id_utilisateur = $1
		`, idUtilisateur).Scan(&user.Nom, &user.Prenom, &user.Email, &user.Telephone, &user.Adresse, &user.Ville, &user.CodePostal, &user.Pays, &nom_entreprise, &type_commerce, &numero_siret)
		if rowsErr != nil {
			fmt.Printf("Erreur : %v", rowsErr)
			http.Error(response, "Erreur récupération valeurs commerçant", http.StatusInternalServerError)
			return
		}

		data := CommercantData{
			Commercant:     user,
			Nom_Entreprise: nom_entreprise,
			Type_Commerce:  type_commerce,
			Numero_Siret:   numero_siret,
		}

		tmpl, errTmpl := template.ParseFiles("./templates/commercant/profil_commercant.html")
		if errTmpl != nil {
			fmt.Printf("Erreur : %v", errTmpl)
			http.Error(response, "Erreur parsefiles html", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "COMMERCANT")
}

func ModifierProfilCommercant(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur récupération session", http.StatusInternalServerError)
			return
		}

		idUtilisateur, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			http.Error(response, "Erreur récupération ID", http.StatusInternalServerError)
			return
		}

		if errForm := request.ParseForm(); errForm != nil {
			fmt.Printf("Erreur : %v", errForm)
			http.Error(response, "Erreur formulaire", http.StatusInternalServerError)
			return
		}

		nom := request.FormValue("nom")
		prenom := request.FormValue("prenom")
		email := request.FormValue("email")
		telephone := request.FormValue("telephone")
		adresse := request.FormValue("adresse")
		ville := request.FormValue("ville")
		codePostal := request.FormValue("code_postal")
		pays := request.FormValue("pays")

		_, errExec := database.Exec(`
				UPDATE utilisateur SET
					nom         = COALESCE(NULLIF($1, ''), nom),
					prenom      = COALESCE(NULLIF($2, ''), prenom),
					email       = COALESCE(NULLIF($3, ''), email),
					telephone   = COALESCE(NULLIF($4, ''), telephone),
					adresse     = COALESCE(NULLIF($5, ''), adresse),
					ville       = COALESCE(NULLIF($6, ''), ville),
					code_postal = COALESCE(NULLIF($7, ''), code_postal),
					pays        = COALESCE(NULLIF($8, ''), pays)
			WHERE id_utilisateur = $9	
		`, nom, prenom, email, telephone, adresse, ville, codePostal, pays, idUtilisateur)
		if errExec != nil {
			fmt.Printf("Erreur : %v", errExec)
			http.Error(response, "Erreur update informations", http.StatusInternalServerError)
			return
		}

		http.Redirect(response, request, "/commercant/profil", http.StatusSeeOther)

	}, "COMMERCANT")
}

func ModificationMotDePasseCommercant(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur récupération session", http.StatusInternalServerError)
			return
		}

		idUtilisateur, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			http.Error(response, "Erreur récupération ID", http.StatusInternalServerError)
			return
		}

		if errForm := request.ParseForm(); errForm != nil {
			fmt.Printf("Erreur : %v", errForm)
			http.Error(response, "Erreur formulaire", http.StatusInternalServerError)
			return
		}

		motDePasseActuelSaisi := request.FormValue("mot_de_passe_actuel")
		nouveauMotDePasse := request.FormValue("nouveau_mot_de_passe")
		confirmation := request.FormValue("confirmation")

		if nouveauMotDePasse != confirmation {
			http.Error(response, "Mot de passe non correspondant", http.StatusBadRequest)
			return
		}

		if nouveauMotDePasse == motDePasseActuelSaisi {
			http.Error(response, "Mot de passe identique", http.StatusBadRequest)
			return
		}

		var hashActuel string
		errSelect := database.QueryRow(`SELECT mot_de_passe FROM utilisateur WHERE id_utilisateur = $1`, idUtilisateur).Scan(&hashActuel)
		if errSelect != nil {
			fmt.Printf("Erreur : %v", errSelect)
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}

		if errCompare := bcrypt.CompareHashAndPassword([]byte(hashActuel), []byte(motDePasseActuelSaisi)); errCompare != nil {
			http.Error(response, "Mot de passe non correspondant", http.StatusBadRequest)
			return
		}

		motDePasseFinal, bcryptErr := bcrypt.GenerateFromPassword([]byte(nouveauMotDePasse), 10)
		if bcryptErr != nil {
			fmt.Printf("Erreur : %v", bcryptErr)
			http.Error(response, "Erreur cryptage mot de passe", http.StatusInternalServerError)
			return
		}

		_, errExec := database.Exec(`UPDATE utilisateur SET mot_de_passe = $1 WHERE id_utilisateur = $2`, motDePasseFinal, idUtilisateur)
		if errExec != nil {
			fmt.Printf("Erreur : %v", errExec)
			http.Error(response, "Erreur update mot de passe", http.StatusInternalServerError)
			return
		}

		http.Redirect(response, request, "/commercant/profil", http.StatusSeeOther)

	}, "COMMERCANT")
}
