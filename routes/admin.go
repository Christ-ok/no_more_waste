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

type AdminDashboardDataBenevole struct {
	Benevoles         []models.BenevoleAffichageDashboard
	BenevoleRecherche string
}

func PageCreerBenevole(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "templates/admin_general/creer_benevole.html")
}

func DashboardAdministrateurBenevoles(database *sql.DB) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur lors de la récupération de la session", http.StatusInternalServerError)
			return
		}

		_, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			http.Error(response, "Utilisateur non identifié", http.StatusUnauthorized)
			return
		}

		rowsBenevole, errBenevole := database.Query(`SELECT u.id_utilisateur, u.nom, u.prenom, u.email, u.telephone, u.adresse, u.ville, u.pays, b.statut FROM utilisateur u 
                                                    JOIN benevole b ON b.id_utilisateur = u.id_utilisateur`,
		)

		if errBenevole != nil {
			fmt.Printf("Voici l'erreur à la ligne 25 : %v\n", errBenevole)
			http.Error(response, "Erreur lors de la récupération des bénévoles", http.StatusInternalServerError)
			return
		}
		defer rowsBenevole.Close()

		var benevoles_List []models.BenevoleAffichageDashboard

		for rowsBenevole.Next() {
			var benevole models.BenevoleAffichageDashboard

			errQuery := rowsBenevole.Scan(&benevole.IDUtilisateur,
				&benevole.Nom,
				&benevole.Prenom,
				&benevole.Email,
				&benevole.Telephone,
				&benevole.Adresse,
				&benevole.Ville,
				&benevole.Pays,
				&benevole.Statut)

			if errQuery != nil {
				fmt.Printf("Erreur lors du Scan : %v\n", errQuery)
				http.Error(response, "Erreur lors de la lecture des données", http.StatusInternalServerError)
				return
			}
			benevoles_List = append(benevoles_List, benevole)
		}

		data := AdminDashboardDataBenevole{
			Benevoles: benevoles_List,
		}

		tmpl, err := template.ParseFiles("./templates/admin_general/benevoles.html")
		if err != nil {
			fmt.Printf("Erreur chargement template : %v\n", err)
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)
	}
}

func CreateBenevole(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		if request.Method != "POST" {
			http.Error(response, "Méthode incorrecte", http.StatusBadRequest)
			return
		}

		err := request.ParseForm()
		if err != nil {
			http.Error(response, "Données de formulaire invalides", http.StatusBadRequest)
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
		motDePasseInitial := request.FormValue("mot_de_passe")
		permis := request.FormValue("permis") == "on"

		if nom == "" || prenom == "" || email == "" || motDePasseInitial == "" {
			http.Error(response, "Champs obligatoires manquants", http.StatusBadRequest)
			return
		}

		hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(motDePasseInitial), 10)
		if errHash != nil {
			http.Error(response, "Erreur lors du hachage du mot de passe", http.StatusInternalServerError)
			return
		}

		var idRole int
		errRole := database.QueryRow("SELECT id_role FROM role WHERE nom = $1", "BENEVOLE").Scan(&idRole)
		if errRole != nil {
			fmt.Printf("Erreur rôle introuvable : %v\n", errRole)
			http.Error(response, "Rôle utilisateur introuvable", http.StatusInternalServerError)
			return
		}

		tx, errTx := database.Begin()
		if errTx != nil {
			http.Error(response, "Erreur d'initialisation de la transaction", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		var idUtilisateur int
		errInsert := tx.QueryRow(
			`INSERT INTO utilisateur (nom, prenom, email, mot_de_passe, telephone, adresse, ville, code_postal, pays, id_role)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
            RETURNING id_utilisateur`,
			nom, prenom, email, string(hashedPassword), telephone, adresse, ville, codePostal, pays, idRole,
		).Scan(&idUtilisateur)
		if errInsert != nil {
			fmt.Printf("Erreur insertion utilisateur : %v\n", errInsert)
			http.Error(response, "Erreur lors de la création de l'utilisateur", http.StatusInternalServerError)
			return
		}

		_, errExec := tx.Exec(
			`INSERT INTO benevole (id_utilisateur, permis) VALUES ($1, $2)`,
			idUtilisateur, permis,
		)
		if errExec != nil {
			fmt.Printf("Erreur insertion bénévole : %v\n", errExec)
			http.Error(response, "Erreur lors de la création du profil bénévole", http.StatusInternalServerError)
			return
		}

		errCommit := tx.Commit()
		if errCommit != nil {
			http.Error(response, "Erreur de validation de la transaction", http.StatusInternalServerError)
			return
		}

		http.Redirect(response, request, "/admin/benevoles", http.StatusSeeOther)

	}, "ADMIN_GENERAL")
}

func FormModifyBenevole(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {
		id := request.URL.Query().Get("id")
		if id == "" {
			http.Error(response, "ID manquant", http.StatusBadRequest)
			return
		}

		var benevole models.BenevoleAffichageDashboard

		errQuery := database.QueryRow(`SELECT 
				u.id_utilisateur,
				u.nom,
				u.prenom,
				u.email,
				COALESCE(u.telephone, ''),
				COALESCE(u.adresse, ''),
				COALESCE(u.ville, ''),
				COALESCE(u.code_postal, ''),
				COALESCE(u.pays, ''),
				b.permis,
				COALESCE(b.disponibilite, ''),
				COALESCE(b.statut, '')
            FROM utilisateur u
            JOIN benevole b
            ON u.id_utilisateur = b.id_utilisateur
            WHERE u.id_utilisateur = $1         
        `, id).Scan(
			&benevole.IDUtilisateur,
			&benevole.Nom,
			&benevole.Prenom,
			&benevole.Email,
			&benevole.Telephone,
			&benevole.Adresse,
			&benevole.Ville,
			&benevole.CodePostal,
			&benevole.Pays,
			&benevole.Permis,
			&benevole.Disponibilite,
			&benevole.Statut,
		)
		if errQuery != nil {
			fmt.Printf("Erreur sélection bénévole : %v\n", errQuery)
			http.Error(response, "Bénévole introuvable", http.StatusNotFound)
			return
		}

		tmpl, errTmpl := template.ParseFiles("./templates/admin_general/modifier_benevole.html")
		if errTmpl != nil {
			http.Error(response, "Erreur de chargement du formulaire", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, benevole)

	}, "ADMIN_GENERAL", "ADMIN_AGENCE")
}

func ModifyBenevole(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {
		id := request.URL.Query().Get("id")
		if id == "" {
			http.Error(response, "ID manquant", http.StatusBadRequest)
			return
		}

		if request.Method != http.MethodPost {
			http.Error(response, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		errForm := request.ParseForm()
		if errForm != nil {
			http.Error(response, "Données de formulaire invalides", http.StatusBadRequest)
			return
		}

		benevole := models.BenevoleAffichageDashboard{
			Nom:           request.FormValue("nom"),
			Prenom:        request.FormValue("prenom"),
			Email:         request.FormValue("email"),
			Telephone:     request.FormValue("telephone"),
			Adresse:       request.FormValue("adresse"),
			Ville:         request.FormValue("ville"),
			CodePostal:    request.FormValue("code_postal"),
			Pays:          request.FormValue("pays"),
			Permis:        request.FormValue("permis") == "on",
			Disponibilite: request.FormValue("disponibilite"),
			Statut:        request.FormValue("statut"),
		}

		tx, err := database.Begin()
		if err != nil {
			http.Error(response, "Erreur de transaction", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		_, err = tx.Exec(`
            UPDATE utilisateur
            SET
                nom = $1,
                prenom = $2,
                email = $3,
                telephone = $4,
                adresse = $5,
                ville = $6,
                code_postal = $7,
                pays = $8
            WHERE id_utilisateur = $9              
        `,
			benevole.Nom,
			benevole.Prenom,
			benevole.Email,
			benevole.Telephone,
			benevole.Adresse,
			benevole.Ville,
			benevole.CodePostal,
			benevole.Pays,
			id,
		)
		if err != nil {
			fmt.Printf("Erreur modification utilisateur : %v\n", err)
			http.Error(response, "Erreur lors de la mise à jour des données utilisateur", http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec(`
            UPDATE benevole
            SET
                permis = $1,
                disponibilite = $2,
                statut = $3
            WHERE id_utilisateur = $4
        `,
			benevole.Permis,
			benevole.Disponibilite,
			benevole.Statut,
			id,
		)
		if err != nil {
			fmt.Printf("Erreur modification bénévole : %v\n", err)
			http.Error(response, "Erreur lors de la mise à jour du profil bénévole", http.StatusInternalServerError)
			return
		}

		err = tx.Commit()
		if err != nil {
			http.Error(response, "Erreur de validation des modifications", http.StatusInternalServerError)
			return
		}

		http.Redirect(response, request, "/admin/benevoles", http.StatusSeeOther)

	}, "ADMIN_GENERAL", "ADMIN_AGENCE")
}

func DeleteBenevole(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {
		id := request.URL.Query().Get("id")

		if id == "" {
			http.Error(response, "ID manquant", http.StatusBadRequest)
			return
		}

		_, err := database.Exec(`DELETE FROM utilisateur WHERE id_utilisateur = $1`, id)
		if err != nil {
			fmt.Printf("Erreur suppression utilisateur : %v\n", err)
			http.Error(response, "Erreur lors de la suppression de l'utilisateur", http.StatusInternalServerError)
			return
		}

		_, err = database.Exec(`DELETE FROM benevole WHERE id_utilisateur = $1`, id)
		if err != nil {
			fmt.Printf("Erreur suppression bénévole : %v\n", err)
			http.Error(response, "Erreur lors de la suppression du bénévole", http.StatusInternalServerError)
			return
		}

		response.WriteHeader(http.StatusOK)
		fmt.Println("Bénévole supprimé, id :", id)

	}, "ADMIN_GENERAL")
}
