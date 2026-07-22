package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"no_more_waste/middleware"
	"no_more_waste/session"
)

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
