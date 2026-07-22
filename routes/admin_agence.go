package routes

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"no_more_waste/middleware"
	"no_more_waste/models"
	"no_more_waste/session"
)

type AdminAgenceDashboardDisponibilite struct {
	Disponibilites []models.BenevoleDisponibilite
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
