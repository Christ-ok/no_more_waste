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

type BenevoleDashboardDisponibilite struct {
	Disponibilites []models.BenevoleDisponibilite
}

func PageDashboardBenevole(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "templates/benevole/accueil_Benevole.html")
}

func PagePlanningBenevole(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "templates/benevole/planning.html")
}

func BenevoleDisponibilite(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			fmt.Printf("Erreur récupération session : %v\n", sessErr)
			http.Error(response, "Erreur lors de la récupération de la session", http.StatusInternalServerError)
			return
		}

		idUtilisateur, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			fmt.Println("Erreur : id_utilisateur absent de la session.")
			http.Error(response, "Utilisateur non connecté", http.StatusUnauthorized)
			return
		}

		fmt.Printf("ID Utilisateur : %d\n", idUtilisateur)

		var idBenevole int

		err := database.QueryRow(
			"SELECT id_benevole FROM benevole WHERE id_utilisateur = $1",
			idUtilisateur,
		).Scan(&idBenevole)

		if err != nil {
			fmt.Printf("Erreur récupération id_benevole : %v\n", err)
			http.Error(response, "Bénévole introuvable", http.StatusInternalServerError)
			return
		}

		fmt.Printf("ID Bénévole : %d\n", idBenevole)

		if request.Method != "POST" {
			fmt.Printf("Méthode reçue : %s\n", request.Method)
			http.Error(response, "Méthode incorrecte", http.StatusBadRequest)
			return
		}

		errForm := request.ParseForm()
		if errForm != nil {
			fmt.Printf("Erreur ParseForm : %v\n", errForm)
			http.Error(response, "Données de formulaire invalides", http.StatusBadRequest)
			return
		}

		date := request.FormValue("date")
		heureDebut := request.FormValue("heure_debut")
		heureFin := request.FormValue("heure_fin")
		statut := request.FormValue("statut")

		fmt.Printf("Date          : %s\n", date)
		fmt.Printf("Heure début   : %s\n", heureDebut)
		fmt.Printf("Heure fin     : %s\n", heureFin)
		fmt.Printf("Statut        : %s\n", statut)

		if date == "" || heureDebut == "" || heureFin == "" || statut == "" {
			fmt.Println("Erreur : un ou plusieurs champs sont vides.")
			http.Error(response, "Champs obligatoires manquants", http.StatusBadRequest)
			return
		}

		fmt.Println("Tentative d'insertion dans la table disponibilite...")

		_, errExec := database.Exec(`
			INSERT INTO disponibilite
			(date_disponibilite, heure_debut, heure_fin, statut, id_benevole)
			VALUES ($1, $2, $3, $4, $5)
		`, date, heureDebut, heureFin, statut, idBenevole)

		if errExec != nil {
			fmt.Printf("%v\n", errExec)
			http.Error(response, "Erreur lors de la création de la disponibilité", http.StatusInternalServerError)
			return
		}

		fmt.Println("Disponibilité ajoutée avec succès !")

		response.WriteHeader(http.StatusOK)

	}, "BENEVOLE")
}

func DashboardBenevoleDisponibiltes(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {
		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur récupération session", http.StatusInternalServerError)
			return
		}

		idUtilisateur, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			http.Error(response, "ID manquant", http.StatusInternalServerError)
			return
		}

		var idBenevole int
		errBenevole := database.QueryRow(`SELECT id_benevole FROM benevole WHERE id_utilisateur = $1`, idUtilisateur).Scan(&idBenevole)
		if errBenevole != nil {
			fmt.Printf("Erreur récupération id_benevole : %v\n", errBenevole)
			http.Error(response, "Bénévole introuvable", http.StatusInternalServerError)
			return
		}

		rowsDisponibilite, errDisponibilties := database.Query(`SELECT id_disponibilite, date_disponibilite, heure_debut, heure_fin, statut FROM disponibilite 
																WHERE id_benevole = $1 ORDER BY date_disponibilite, heure_debut`,
			idBenevole)

		if errDisponibilties != nil {
			fmt.Printf("Erreur : %v", errDisponibilties)
			http.Error(response, "Erreur récupération disponibilite", http.StatusInternalServerError)
			return
		}
		defer rowsDisponibilite.Close()

		var disponibilite_List []models.BenevoleDisponibilite

		for rowsDisponibilite.Next() {
			var disponibilite models.BenevoleDisponibilite

			errQuery := rowsDisponibilite.Scan(
				&disponibilite.ID_Disponibilite,
				&disponibilite.Date_Disponibilite,
				&disponibilite.Heure_Debut,
				&disponibilite.Heure_Fin,
				&disponibilite.Statut,
			)
			if errQuery != nil {
				fmt.Printf("Erreur : %v", errQuery)
				http.Error(response, "Erreur scan disponibilite", http.StatusInternalServerError)
				return
			}

			disponibilite_List = append(disponibilite_List, disponibilite)
		}

		if errRows := rowsDisponibilite.Err(); errRows != nil {
			fmt.Printf("Erreur itération disponibilités : %v\n", errRows)
			http.Error(response, "Erreur lecture disponibilite", http.StatusInternalServerError)
			return
		}

		data := BenevoleDashboardDisponibilite{
			Disponibilites: disponibilite_List,
		}

		tmpl, errTmpl := template.ParseFiles("./templates/benevole/disponibilites.html")
		if errTmpl != nil {
			fmt.Printf("Erreur : %v", errTmpl)
			http.Error(response, "Erreur template disponibilite", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "BENEVOLE")
}

func ModifierDisponibilite(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur lors de la récupération de la session", http.StatusInternalServerError)
			return
		}

		idUtilisateur, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			http.Error(response, "Utilisateur non identifié", http.StatusUnauthorized)
			return
		}

		var idBenevole int
		errBenevole := database.QueryRow(`SELECT id_benevole FROM benevole WHERE id_utilisateur = $1`, idUtilisateur).Scan(&idBenevole)
		if errBenevole != nil {
			fmt.Printf("Erreur : %v", errBenevole)
			http.Error(response, "Erreur récupération bénévole", http.StatusInternalServerError)
			return
		}

		errForm := request.ParseForm()
		if errForm != nil {
			http.Error(response, "Données de formulaire invalides", http.StatusBadRequest)
			return
		}

		idDisponibilite := request.FormValue("id_disponibilite")
		date := request.FormValue("date")
		heure_debut := request.FormValue("heure_debut")
		heure_fin := request.FormValue("heure_fin")
		statut := request.FormValue("statut")

		if idDisponibilite == "" || date == "" || heure_debut == "" || heure_fin == "" || statut == "" {
			http.Error(response, "Champs obligatoires manquants", http.StatusBadRequest)
			return
		}

		tx, errTx := database.Begin()
		if errTx != nil {
			fmt.Printf("Erreur : %v", errTx)
			http.Error(response, "Erreur création transaction", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		result, errExec := tx.Exec(`UPDATE disponibilite SET date_disponibilite = $1, heure_debut = $2, heure_fin = $3, statut = $4 
									WHERE id_benevole = $5 AND id_disponibilite = $6`, date, heure_debut, heure_fin, statut, idBenevole, idDisponibilite)
		if errExec != nil {
			fmt.Printf("Erreur : %v", errExec)
			http.Error(response, "Erreur modification disponibilité", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(response, "Aucune ligne affectée par la transaction", http.StatusForbidden)
			return
		}

		_, errPlanning := tx.Exec(`UPDATE planning SET date = $1, heure_debut = $2, heure_fin = $3
									WHERE id_disponibilite = $4`, date, heure_debut, heure_fin, idDisponibilite)
		if errPlanning != nil {
			fmt.Printf("Erreur : %v", errPlanning)
			http.Error(response, "Erreur Modification planning", http.StatusInternalServerError)
			return
		}

		errCommit := tx.Commit()
		if errCommit != nil {
			fmt.Printf("Erreur : %v", errCommit)
			http.Error(response, "Erreur commit de la transaction", http.StatusInternalServerError)
			return
		}

		fmt.Println("Disponibilité mise à jour")
		response.WriteHeader(http.StatusOK)

	}, "BENEVOLE")
}

func DeleteDisponibilite(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur récupération session", http.StatusInternalServerError)
			return
		}

		idUtilisateur, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			http.Error(response, "Erreur récupération id", http.StatusInternalServerError)
			return
		}

		var idBenevole int
		errQuery := database.QueryRow(`SELECT id_benevole FROM benevole WHERE id_utilisateur = $1`, idUtilisateur).Scan(&idBenevole)
		if errQuery != nil {
			fmt.Printf("Erreur : %v", errQuery)
			http.Error(response, "Erreur récupération ID", http.StatusInternalServerError)
			return
		}

		errForm := request.ParseForm()
		if errForm != nil {
			http.Error(response, "Données de formulaire invalides", http.StatusBadRequest)
			return
		}

		idDisponibilite := request.FormValue("id_disponibilite")
		if idDisponibilite == "" {
			http.Error(response, "ID manquant", http.StatusBadRequest)
			return
		}

		tx, errTx := database.Begin()
		if errTx != nil {
			fmt.Printf("Erreur : %v", errTx)
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		_, errDeletePlanning := tx.Exec(`DELETE FROM planning WHERE id_disponibilite = $1`, idDisponibilite)
		if errDeletePlanning != nil {
			fmt.Printf("Erreur suppression planning associé : %v\n", errDeletePlanning)
			http.Error(response, "Erreur suppression du planning associé", http.StatusInternalServerError)
			return
		}

		result, errExec := tx.Exec(`DELETE FROM disponibilite WHERE id_disponibilite = $1 AND id_benevole = $2`, idDisponibilite, idBenevole)
		if errExec != nil {
			fmt.Printf("Erreur : %v", errExec)
			http.Error(response, "Erreur suppression disponibilité", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(response, "Erreur affectation de lignes", http.StatusInternalServerError)
			return
		}

		errCommit := tx.Commit()
		if errCommit != nil {
			fmt.Printf("Erreur : %v", errCommit)
			http.Error(response, "Erreur commit transaction", http.StatusInternalServerError)
			return
		}

		fmt.Println("Disponibilité supprimé avec succès")
		response.WriteHeader(http.StatusOK)

	}, "BENEVOLE")
}
