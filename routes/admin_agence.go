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
)

type AdminAgenceDashboardDisponibilite struct {
	Disponibilites []models.BenevoleDisponibilite
}

type AdminAgenceDashboardPlanning struct {
	Planning models.Planning
}

func DashboardAdminAgenceBenevoles(database *sql.DB) http.HandlerFunc {
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

		rowsBenevole, errBenevole := database.Query(`SELECT u.id_utilisateur, u.nom, u.prenom, u.email, u.adresse, u.ville, u.pays FROM utilisateur u
													JOIN benevole b ON b.id_utilisateur = u.id_utilisateur

		`)

		if errBenevole != nil {
			fmt.Printf("Voici l'erreur à la ligne 32 : %v\n", errBenevole)
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
				&benevole.Adresse,
				&benevole.Ville,
				&benevole.Pays,
			)

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

		tmpl, errTmpl := template.ParseFiles("./templates/admin_agence/benevoles.html")
		if errTmpl != nil {
			fmt.Printf("Erreur chargement template : %v\n", errTmpl)
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)
	}
}

func DashboardAdminAgenceGererDisponibilite(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {
		id := request.URL.Query().Get("id")
		if id == "" {
			http.Error(response, "ID manquant", http.StatusBadRequest)
			return
		}

		var idBenevole int
		err := database.QueryRow(`SELECT id_benevole FROM benevole WHERE id_utilisateur = $1`, id).Scan(&idBenevole)
		if err != nil {
			fmt.Printf("Voici l'erreur à la ligne 93 : %v\n", err)
			http.Error(response, "Erreur lors de la récupération du nom bénévole", http.StatusInternalServerError)
			return
		}

		rowsDisponibilite, errQuery := database.Query(`SELECT d.id_disponibilite, u.nom, u.prenom, d.date_disponibilite, d.heure_debut, d.heure_fin, d.statut FROM disponibilite d
													JOIN benevole b ON b.id_benevole = d.id_benevole
													LEFT JOIN utilisateur u ON u.id_utilisateur = b.id_utilisateur												
													WHERE d.id_benevole = $1`, idBenevole)
		if errQuery != nil {
			fmt.Printf("Voici l'erreur à la ligne 101 : %v\n", errQuery)
			http.Error(response, "Erreur lors de la récupération des disponibilités du bénévole", http.StatusInternalServerError)
			return
		}
		defer rowsDisponibilite.Close()

		var disponibilites_List []models.BenevoleDisponibilite

		for rowsDisponibilite.Next() {
			var disponibilite models.BenevoleDisponibilite

			errScan := rowsDisponibilite.Scan(&disponibilite.ID_Disponibilite,
				&disponibilite.Nom,
				&disponibilite.Prenom,
				&disponibilite.Date_Disponibilite,
				&disponibilite.Heure_Debut,
				&disponibilite.Heure_Fin,
				&disponibilite.Statut,
			)
			if errScan != nil {
				fmt.Println("L'erreur est à la ligne 117 : %v", errScan)
				http.Error(response, errScan.Error(), http.StatusInternalServerError)
				return
			}

			disponibilites_List = append(disponibilites_List, disponibilite)
		}

		data := AdminAgenceDashboardDisponibilite{
			Disponibilites: disponibilites_List,
		}

		tmpl, errTmpl := template.ParseFiles("./templates/admin_agence/disponibilites_benevoles.html")
		if errTmpl != nil {
			fmt.Printf("Erreur chargement template : %v\n", errTmpl)
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)
	}, "ADMIN_AGENCE")
}

func FormCreatePlanning(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {
		idDisponibilite := request.URL.Query().Get("id")
		if idDisponibilite == "" {
			http.Error(response, "ID manquant", http.StatusBadRequest)
			return
		}

		errForm := request.ParseForm()
		if errForm != nil {
			http.Error(response, "Données de formulaire invalides", http.StatusBadRequest)
			return
		}

		var planning models.Planning

		errExec := database.QueryRow("SELECT id_benevole, id_disponibilite, date_disponibilite, heure_debut, heure_fin, statut FROM disponibilite WHERE id_disponibilite = $1", idDisponibilite).Scan(&planning.ID_Benevole, &planning.ID_Disponibilite, &planning.Date, &planning.Heure_Debut, &planning.Heure_Fin, &planning.Statut)
		if errExec != nil {
			fmt.Printf("Erreur : %v", errExec)
			http.Error(response, "Erreur select dans la base de données", http.StatusInternalServerError)
			return
		}

		data := AdminAgenceDashboardPlanning{
			Planning: planning,
		}

		fmt.Println("ID disponibilité :", planning.ID_Disponibilite)
		fmt.Println("ID bénévole :", planning.ID_Benevole)
		fmt.Println("Date :", planning.Date)
		fmt.Println("Début :", planning.Heure_Debut)
		fmt.Println("Fin :", planning.Heure_Fin)
		fmt.Println("Statut :", planning.Statut)

		tmpl, errTmpl := template.ParseFiles("./templates/admin_agence/creer_planning.html")
		if errTmpl != nil {
			fmt.Printf("Erreur chargement template : %v\n", errTmpl)
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "ADMIN_AGENCE")
}

func DashboardAdminAgenceCreerPlanning(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {
		errForm := request.ParseForm()
		if errForm != nil {
			http.Error(response, "Données de formulaire invalides", http.StatusBadRequest)
			return
		}

		idDisponibilite, errorDisponibilite := strconv.Atoi(request.FormValue("id_disponibilite"))
		if errorDisponibilite != nil {
			fmt.Printf("Erreur : %v", errorDisponibilite)
			http.Error(response, "Erreur conversion idDisponibilite", http.StatusInternalServerError)
			return
		}

		idBenevole, errorBenevole := strconv.Atoi(request.FormValue("id_benevole"))
		if errorBenevole != nil {
			fmt.Printf("Erreur : %v", errorBenevole)
			http.Error(response, "Erreur conversion idBenevole", http.StatusInternalServerError)
			return
		}

		date := request.FormValue("date")
		heure_debut := request.FormValue("heure_debut")
		heure_fin := request.FormValue("heure_fin")
		statut := request.FormValue("statut")

		if date == "" || heure_debut == "" || heure_fin == "" {
			http.Error(response, "Champs obligatoires manquants", http.StatusBadRequest)
			return
		}

		_, errExec := database.Exec("INSERT INTO planning (id_benevole, id_disponibilite, date, heure_debut, heure_fin, statut) VALUES ($1, $2, $3, $4, $5, $6)", idBenevole, idDisponibilite, date, heure_debut, heure_fin, statut)
		if errExec != nil {
			fmt.Printf("Erreur insertion planning : %v\n", errExec)
			http.Error(response, "Erreur lors de la création du planning", http.StatusInternalServerError)
			return
		}

		_, err := database.Exec(`UPDATE disponibilite SET statut ='RESERVE' WHERE id_disponibilite = $1`, idDisponibilite)
		if err != nil {
			fmt.Printf("Erreur : %v", err)
			http.Error(response, "Erreur update disponibilite", http.StatusInternalServerError)
			return
		}

		fmt.Println("Disponibilité créée avec succès !")
		http.Redirect(response, request, "/admin-agence/benevoles/disponibilites", http.StatusSeeOther)

	}, "ADMIN_AGENCE")
}
