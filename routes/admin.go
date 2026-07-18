package routes

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"no_more_waste/middleware"
	"no_more_waste/models"
	"no_more_waste/session"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type AdminDashboardDataBenevole struct {
	Benevoles         []models.BenevoleAffichageDashboard
	BenevoleRecherche string
}

type AdminDashboardDataCommercant struct {
	Commercants        []models.CommercantAffichageDashboard
	CommercantRecherce string
}

type AdminDashboardDataAdministrateur struct {
	Administrateurs         []models.AdministrateurAffichageDashboard
	AdministrateurRecherche string
	Agences                 []models.Agence
}

type AdminModifyDataAdministrateur struct {
	Administrateur models.AdministrateurAffichageDashboard
	Agences        []models.Agence
}

func PageCreerBenevole(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "templates/admin_general/creer_benevole.html")
}

func PageCreerCommercant(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "templates/admin_general/creer_commercant.html")
}

func getAgences(database *sql.DB) ([]models.Agence, error) {
	rows, err := database.Query(`SELECT id_agence, nom, ville FROM agence`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agences []models.Agence
	for rows.Next() {
		var agence models.Agence
		if err := rows.Scan(&agence.IDAgence, &agence.Nom, &agence.Ville); err != nil {
			return nil, err
		}
		agences = append(agences, agence)
	}
	return agences, nil
}

func PageCreerAdministrateur(database *sql.DB) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		agences, err := getAgences(database)
		if err != nil {
			fmt.Printf("Erreur récupération agences : %v\n", err)
			http.Error(response, "Erreur lors de la récupération des agences", http.StatusInternalServerError)
			return
		}

		data := AdminDashboardDataAdministrateur{
			Agences: agences,
		}

		tmpl, errTmpl := template.ParseFiles("./templates/admin_general/creer_administrateur.html")
		if errTmpl != nil {
			fmt.Printf("Erreur chargement template : %v\n", errTmpl)
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)
	}
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

		_, err := database.Exec(`DELETE FROM benevole WHERE id_utilisateur = $1`, id)
		if err != nil {
			fmt.Printf("Erreur suppression bénévole : %v\n", err)
			http.Error(response, "Erreur lors de la suppression du bénévole", http.StatusInternalServerError)
			return
		}

		_, err = database.Exec(`DELETE FROM utilisateur WHERE id_utilisateur = $1`, id)
		if err != nil {
			fmt.Printf("Erreur suppression utilisateur : %v\n", err)
			http.Error(response, "Erreur lors de la suppression de l'utilisateur", http.StatusInternalServerError)
			return
		}

		response.WriteHeader(http.StatusOK)
		fmt.Println("Bénévole supprimé, id :", id)

	}, "ADMIN_GENERAL")
}

func DashboardAdministrateurCommercant(database *sql.DB) http.HandlerFunc {
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

		rowsCommercant, errCommercant := database.Query(`SELECT u.id_utilisateur, u.nom, u.prenom, u.email, u.telephone, u.adresse, u.ville, u.pays, c.nom_entreprise, c.type_commerce, c.numero_siret, c.statut
															FROM utilisateur u JOIN commercant c ON u.id_utilisateur = c.id_utilisateur`,
		)
		if errCommercant != nil {
			fmt.Printf("Voici l'erreur à la ligne 25 : %v\n", errCommercant)
			http.Error(response, "Erreur lors de la récupération des bénévoles", http.StatusInternalServerError)
			return
		}
		defer rowsCommercant.Close()

		var commercants_List []models.CommercantAffichageDashboard

		for rowsCommercant.Next() {
			var commercant models.CommercantAffichageDashboard

			errQuery := rowsCommercant.Scan(&commercant.IDutilisateur,
				&commercant.Nom,
				&commercant.Prenom,
				&commercant.Email,
				&commercant.Telephone,
				&commercant.Adresse,
				&commercant.Ville,
				&commercant.Pays,
				&commercant.NomEntriprise,
				&commercant.TypeCommerce,
				&commercant.NumeroSiret,
				&commercant.Statut,
			)
			if errQuery != nil {
				fmt.Printf("Erreur lors du Scan : %v\n", errQuery)
				http.Error(response, "Erreur lors de la lecture des données", http.StatusInternalServerError)
				return
			}
			commercants_List = append(commercants_List, commercant)
		}

		data := AdminDashboardDataCommercant{
			Commercants: commercants_List,
		}

		tmpl, errTmpl := template.ParseFiles("./templates/admin_general/commercants.html")
		if errTmpl != nil {
			fmt.Printf("Erreur chargement template : %v\n", errTmpl)
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)
	}
}

func CreateCommercant(database *sql.DB) http.HandlerFunc {
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
		nomEntriprise := request.FormValue("nom_entreprise")
		typeCommerce := request.FormValue("type_commerce")
		numeroSiret := request.FormValue("numero_siret")
		statut := request.FormValue("statut")

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
		errRole := database.QueryRow("SELECT id_role FROM role WHERE nom = $1", "COMMERCANT").Scan(&idRole)
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
			`INSERT INTO commercant (id_utilisateur, nom_entreprise, type_commerce, numero_siret, statut)
			VALUES ($1, $2, $3, $4, $5)`, idUtilisateur, nomEntriprise, typeCommerce, numeroSiret, statut,
		)
		if errExec != nil {
			fmt.Printf("Erreur insertion commercant : %v\n", errExec)
			http.Error(response, "Erreur lors de la création du profil commercant", http.StatusInternalServerError)
			return
		}

		errCommit := tx.Commit()
		if errCommit != nil {
			http.Error(response, "Erreur de validation de la transaction", http.StatusInternalServerError)
			return
		}

		http.Redirect(response, request, "/admin/commercants", http.StatusSeeOther)

	}, "ADMIN_GENERAL")
}

func FormModifyCommercant(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {
		id := request.URL.Query().Get("id")
		if id == "" {
			http.Error(response, "ID manquant", http.StatusBadRequest)
			return
		}

		var commercant models.CommercantAffichageDashboard

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
				COALESCE(c.nom_entreprise, ''),
				COALESCE(c.type_commerce, ''),
				COALESCE(c.numero_siret, ''),
				COALESCE(c.statut, '')
            FROM utilisateur u
            JOIN commercant c
            ON u.id_utilisateur = c.id_utilisateur
            WHERE u.id_utilisateur = $1         
        `, id).Scan(
			&commercant.IDutilisateur,
			&commercant.Nom,
			&commercant.Prenom,
			&commercant.Email,
			&commercant.Telephone,
			&commercant.Adresse,
			&commercant.Ville,
			&commercant.CodePostal,
			&commercant.Pays,
			&commercant.NomEntriprise,
			&commercant.TypeCommerce,
			&commercant.NumeroSiret,
			&commercant.Statut,
		)
		if errQuery != nil {
			fmt.Printf("Erreur sélection commercant : %v\n", errQuery)
			http.Error(response, "Commercant introuvable", http.StatusNotFound)
			return
		}

		tmpl, errTmpl := template.ParseFiles("./templates/admin_general/modifier_commercant.html")
		if errTmpl != nil {
			http.Error(response, "Erreur de chargement du formulaire", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, commercant)

	}, "ADMIN_GENERAL", "ADMIN_AGENCE")
}

func ModifyCommercant(database *sql.DB) http.HandlerFunc {
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

		commercant := models.CommercantAffichageDashboard{
			Nom:           request.FormValue("nom"),
			Prenom:        request.FormValue("prenom"),
			Email:         request.FormValue("email"),
			Telephone:     request.FormValue("telephone"),
			Adresse:       request.FormValue("adresse"),
			Ville:         request.FormValue("ville"),
			CodePostal:    request.FormValue("code_postal"),
			Pays:          request.FormValue("pays"),
			NomEntriprise: request.FormValue("nom_entreprise"),
			TypeCommerce:  request.FormValue("type_commerce"),
			NumeroSiret:   request.FormValue("numero_siret"),
			Statut:        request.FormValue("statut"),
		}

		tx, err := database.Begin()
		if err != nil {
			http.Error(response, "Erreur de transaction", http.
				StatusInternalServerError)
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
			commercant.Nom,
			commercant.Prenom,
			commercant.Email,
			commercant.Telephone,
			commercant.Adresse,
			commercant.Ville,
			commercant.CodePostal,
			commercant.Pays,
			id,
		)
		if err != nil {
			fmt.Printf("Erreur modification utilisateur : %v\n", err)
			http.Error(response, "Erreur lors de la mise à jour des données utilisateur", http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec(`
			UPDATE commercant
			SET
				nom_entreprise = $1,
				type_commerce = $2,
				numero_siret = $3,
				statut = $4
			WHERE id_utilisateur = $5 
		`,
			commercant.NomEntriprise,
			commercant.TypeCommerce,
			commercant.NumeroSiret,
			commercant.Statut,
			id,
		)
		if err != nil {
			fmt.Printf("Erreur modification commercant : %v\n", err)
			http.Error(response, "Erreur lors de la mise à jour du profil commercant", http.StatusInternalServerError)
			return
		}

		err = tx.Commit()
		if err != nil {
			http.Error(response, "Erreur de validation des modifications", http.StatusInternalServerError)
			return
		}

		http.Redirect(response, request, "/admin/commercants", http.StatusSeeOther)

	}, "ADMIN_GENERAL", "ADMIN_AGENCE")
}

func DeleteCommercant(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {
		id := request.URL.Query().Get("id")

		if id == "" {
			http.Error(response, "ID manquant", http.StatusBadRequest)
			return
		}

		_, err := database.Exec(`DELETE FROM commercant WHERE id_utilisateur = $1`, id)
		if err != nil {
			fmt.Printf("Erreur suppression commercant : %v\n", err)
			http.Error(response, "Erreur lors de la suppression du commercant", http.StatusInternalServerError)
			return
		}

		_, err = database.Exec(`DELETE FROM utilisateur WHERE id_utilisateur = $1`, id)
		if err != nil {
			fmt.Printf("Erreur suppression utilisateur : %v\n", err)
			http.Error(response, "Erreur lors de la suppression de l'utilisateur", http.StatusInternalServerError)
			return
		}

		response.WriteHeader(http.StatusOK)
		fmt.Println("Commercant supprimé, id :", id)

	}, "ADMIN_GENERAL")
}

func DashboardAdministrateurAdmins(database *sql.DB) http.HandlerFunc {
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

		rowsAdmin, errAdmin := database.Query(`SELECT u.id_utilisateur, u.nom, u.prenom, u.email, u.telephone, u.adresse, u.ville, u.pays, 
													COALESCE(u.id_agence, 0),
													COALESCE(a.nom, 'Aucune agence'),
													u.etat_compte
												FROM utilisateur u
												JOIN role r ON r.id_role = u.id_role
												LEFT JOIN agence a ON a.id_agence = u.id_agence
												WHERE r.nom = 'ADMIN_AGENCE'`,
		)

		if errAdmin != nil {
			fmt.Printf("Erreur récupération administrateurs : %v\n", errAdmin)
			http.Error(response, "Erreur lors de la récupération des administrateurs", http.StatusInternalServerError)
			return
		}
		defer rowsAdmin.Close()

		var administrateurs_List []models.AdministrateurAffichageDashboard

		for rowsAdmin.Next() {
			var administrateur models.AdministrateurAffichageDashboard

			errQuery := rowsAdmin.Scan(
				&administrateur.IDUtilisateur,
				&administrateur.Nom,
				&administrateur.Prenom,
				&administrateur.Email,
				&administrateur.Telephone,
				&administrateur.Adresse,
				&administrateur.Ville,
				&administrateur.Pays,
				&administrateur.IDAgence,
				&administrateur.NomAgence,
				&administrateur.EtatCompte,
			)
			if errQuery != nil {
				fmt.Printf("Erreur lors du Scan : %v\n", errQuery)
				http.Error(response, "Erreur lors de la lecture des données", http.StatusInternalServerError)
				return
			}
			administrateurs_List = append(administrateurs_List, administrateur)
		}

		data := AdminDashboardDataAdministrateur{
			Administrateurs: administrateurs_List,
		}

		tmpl, err := template.ParseFiles("./templates/admin_general/administrateurs.html")
		if err != nil {
			fmt.Printf("Erreur chargement template : %v\n", err)
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)
	}
}

func CreateAdministrateur(database *sql.DB) http.HandlerFunc {
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
		idAgenceString := request.FormValue("id_agence")

		if nom == "" || prenom == "" || email == "" || motDePasseInitial == "" {
			http.Error(response, "Champs obligatoires manquants", http.StatusBadRequest)
			return
		}

		idAgence, errConv := strconv.Atoi(idAgenceString)
		if errConv != nil {
			http.Error(response, "Agence invalide", http.StatusBadRequest)
			return
		}

		var idAgenceSQL sql.NullInt64
		if idAgence > 0 {
			idAgenceSQL = sql.NullInt64{Int64: int64(idAgence), Valid: true}
		} else {
			idAgenceSQL = sql.NullInt64{Valid: false}
		}

		hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(motDePasseInitial), 10)
		if errHash != nil {
			http.Error(response, "Erreur lors du hachage du mot de passe", http.StatusInternalServerError)
			return
		}

		var idRole int
		errRole := database.QueryRow("SELECT id_role FROM role WHERE nom = $1", "ADMIN_AGENCE").Scan(&idRole)
		if errRole != nil {
			fmt.Printf("Erreur rôle introuvable : %v\n", errRole)
			http.Error(response, "Rôle utilisateur introuvable", http.StatusInternalServerError)
			return
		}

		_, errInsert := database.Exec(`INSERT INTO utilisateur (nom, prenom, email, mot_de_passe, telephone, adresse, ville, code_postal, pays, id_role, id_agence)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			nom, prenom, email, string(hashedPassword), telephone, adresse, ville, codePostal, pays, idRole, idAgenceSQL,
		)
		if errInsert != nil {
			fmt.Printf("Erreur insertion administrateur : %v\n", errInsert)
			http.Error(response, "Erreur lors de la création de l'administrateur", http.StatusInternalServerError)
			return
		}

		http.Redirect(response, request, "/admin/administrateurs", http.StatusSeeOther)
	}, "ADMIN_GENERAL")
}

func FormModifyAdministrateur(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {
		id := request.URL.Query().Get("id")
		if id == "" {
			http.Error(response, "ID manquant", http.StatusBadRequest)
			return
		}

		var administrateur models.AdministrateurAffichageDashboard

		errQuery := database.QueryRow(`
			SELECT u.id_utilisateur, u.nom, u.prenom, u.email,
			       COALESCE(u.telephone, ''), COALESCE(u.adresse, ''), COALESCE(u.ville, ''),
			       COALESCE(u.code_postal, ''), COALESCE(u.pays, ''), COALESCE(u.id_agence, 0)
			FROM utilisateur u
			WHERE u.id_utilisateur = $1
		`, id).Scan(
			&administrateur.IDUtilisateur,
			&administrateur.Nom,
			&administrateur.Prenom,
			&administrateur.Email,
			&administrateur.Telephone,
			&administrateur.Adresse,
			&administrateur.Ville,
			&administrateur.CodePostal,
			&administrateur.Pays,
			&administrateur.IDAgence,
		)
		if errQuery != nil {
			fmt.Printf("Erreur sélection administrateur : %v\n", errQuery)
			http.Error(response, "Administrateur introuvable", http.StatusNotFound)
			return
		}

		agences, errAgences := getAgences(database)
		if errAgences != nil {
			fmt.Printf("Erreur récupération agences : %v\n", errAgences)
			http.Error(response, "Erreur lors de la récupération des agences", http.StatusInternalServerError)
			return
		}

		data := AdminModifyDataAdministrateur{
			Administrateur: administrateur,
			Agences:        agences,
		}

		tmpl, errTmpl := template.ParseFiles("./templates/admin_general/modifier_administrateur.html")
		if errTmpl != nil {
			http.Error(response, "Erreur de chargement du formulaire", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "ADMIN_GENERAL")
}

func ModifyAdministrateur(database *sql.DB) http.HandlerFunc {
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

		idAgenceStr := request.FormValue("id_agence")
		idAgence, errConv := strconv.Atoi(idAgenceStr)
		if errConv != nil {
			http.Error(response, "Agence invalide", http.StatusBadRequest)
			return
		}

		_, err := database.Exec(`
			UPDATE utilisateur
			SET
				nom = $1,
				prenom = $2,
				email = $3,
				telephone = $4,
				adresse = $5,
				ville = $6,
				code_postal = $7,
				pays = $8,
				id_agence = $9
			WHERE id_utilisateur = $10
		`,
			request.FormValue("nom"),
			request.FormValue("prenom"),
			request.FormValue("email"),
			request.FormValue("telephone"),
			request.FormValue("adresse"),
			request.FormValue("ville"),
			request.FormValue("code_postal"),
			request.FormValue("pays"),
			idAgence,
			id,
		)
		if err != nil {
			fmt.Printf("Erreur modification administrateur : %v\n", err)
			http.Error(response, "Erreur lors de la mise à jour de l'administrateur", http.StatusInternalServerError)
			return
		}

		http.Redirect(response, request, "/admin/administrateurs", http.StatusSeeOther)

	}, "ADMIN_GENERAL")
}

func DeleteAdministrateur(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {
		id := request.URL.Query().Get("id")

		if id == "" {
			http.Error(response, "ID manquant", http.StatusBadRequest)
			return
		}

		_, err := database.Exec(`DELETE FROM utilisateur WHERE id_utilisateur = $1`, id)
		if err != nil {
			fmt.Printf("Erreur suppression administrateur : %v\n", err)
			http.Error(response, "Erreur lors de la suppression de l'administrateur", http.StatusInternalServerError)
			return
		}

		response.WriteHeader(http.StatusOK)
		fmt.Println("Administrateur supprimé, id :", id)

	}, "ADMIN_GENERAL")
}
