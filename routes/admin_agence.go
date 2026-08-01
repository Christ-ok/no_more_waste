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

type AdminAgenceDashboardPlanningAfficher struct {
	Planning []models.PlanningAfficheDashboard
}

type AdminAgenceDashboardServicesAfficher struct {
	Services    []models.Service
	Competences []models.Competence
}
type ServiceCreerData struct {
	Competences []models.Competence
}

func DashboardAdminAgenceBenevoles(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur lors de la récupération de la session", http.StatusInternalServerError)
			return
		}

		idUtilisateurAdmin, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			http.Error(response, "Utilisateur non identifié", http.StatusUnauthorized)
			return
		}

		idAgence, errAgence := middleware.GetIDAgenceUtilisateur(database, idUtilisateurAdmin)
		if errAgence != nil {
			fmt.Printf("Erreur récupération agence : %v\n", errAgence)
			http.Error(response, "Agence introuvable pour cet utilisateur", http.StatusForbidden)
			return
		}

		rowsBenevole, errBenevole := database.Query(`SELECT u.id_utilisateur, u.nom, u.prenom, u.email, u.adresse, u.ville, u.pays FROM utilisateur u
													JOIN benevole b ON b.id_utilisateur = u.id_utilisateur
													WHERE u.id_agence = $1

		`, idAgence)

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

	}, "ADMIN_AGENCE")
}

func DashboardAdminAgenceGererDisponibilite(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {
		id := request.URL.Query().Get("id")
		if id == "" {
			http.Error(response, "ID manquant", http.StatusBadRequest)
			return
		}

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur lors de la récupération de la session", http.StatusInternalServerError)
			return
		}

		idUtilisateurAdmin, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			http.Error(response, "Utilisateur non identifié", http.StatusUnauthorized)
			return
		}

		idAgence, errAgence := middleware.GetIDAgenceUtilisateur(database, idUtilisateurAdmin)
		if errAgence != nil {
			http.Error(response, "Agence introuvable pour cet utilisateur", http.StatusForbidden)
			return
		}

		var idBenevole int
		err := database.QueryRow(`SELECT b.id_benevole FROM benevole b 
								JOIN utilisateur u ON u.id_utilisateur = b.id_utilisateur
								WHERE b.id_utilisateur = $1 AND u.id_agence = $2`,
			id, idAgence).Scan(&idBenevole)

		if err == sql.ErrNoRows {
			http.Error(response, "Bénévole introuvable dans votre agence", http.StatusForbidden)
			return
		}

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
			fmt.Printf("Voici l'erreur à la ligne 143 : %v\n", errQuery)
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
				fmt.Printf("L'erreur est : %v", errScan)
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

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur lors de la récupération de la session", http.StatusInternalServerError)
			return
		}

		idUtilisateurAdmin, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			http.Error(response, "Utilisateur non identifié", http.StatusUnauthorized)
			return
		}

		idAgence, errAgence := middleware.GetIDAgenceUtilisateur(database, idUtilisateurAdmin)
		if errAgence != nil {
			http.Error(response, "Agence introuvable pour cet utilisateur", http.StatusForbidden)
			return
		}

		var planning models.Planning

		errExec := database.QueryRow(`SELECT d.id_benevole, d.id_disponibilite, d.date_disponibilite, d.heure_debut, d.heure_fin, d.statut FROM disponibilite d
									JOIN benevole b ON b.id_benevole = d.id_benevole
									JOIN utilisateur u ON u.id_utilisateur = b.id_utilisateur 
									WHERE d.id_disponibilite = $1 AND u.id_agence = $2`, idDisponibilite, idAgence).Scan(&planning.ID_Benevole, &planning.ID_Disponibilite, &planning.Date, &planning.Heure_Debut, &planning.Heure_Fin, &planning.Statut)

		if errExec == sql.ErrNoRows {
			http.Error(response, "Disponibilité introuvable dans votre agence", http.StatusForbidden)
			return
		}

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

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur lors de la récupération de la session", http.StatusInternalServerError)
			return
		}

		idUtilisateurAdmin, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			http.Error(response, "Utilisateur non identifié", http.StatusUnauthorized)
			return
		}

		idAgence, errAgence := middleware.GetIDAgenceUtilisateur(database, idUtilisateurAdmin)
		if errAgence != nil {
			http.Error(response, "Agence introuvable pour cet utilisateur", http.StatusForbidden)
			return
		}

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

		var check int
		errCheck := database.QueryRow(`SELECT d.id_disponibilite FROM disponibilite d
										JOIN benevole b ON b.id_benevole = d.id_benevole
										JOIN utilisateur u ON u.id_utilisateur = b.id_utilisateur
										WHERE d.id_disponibilite = $1 AND d.id_benevole = $2 AND u.id_agence = $3 AND d.statut = 'DISPONIBLE'
		`, idDisponibilite, idBenevole, idAgence).Scan(&check)

		if errCheck == sql.ErrNoRows {
			http.Error(response, "Cette disponibilité n'existe pas, n'appartient pas à votre agence, ou a déjà été attribuée", http.StatusForbidden)
			return
		}

		if errCheck != nil {
			fmt.Printf("Erreur vérification : %v\n", errCheck)
			http.Error(response, "Erreur lors de la vérification", http.StatusInternalServerError)
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

func MiseAJourPlanningsExpires(database *sql.DB) error {
	_, err := database.Exec(`
		UPDATE planning SET statut = 'TERMINE'
		WHERE (date + heure_fin) < NOW() 
			AND statut = 'PLANIFIE'
	`)

	return err
}

func DashboardAdminAgenceAfficherPlanning(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur lors de la récupération de la session", http.StatusInternalServerError)
			return
		}

		idUtilisateurAdmin, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			http.Error(response, "Utilisateur non identifié", http.StatusUnauthorized)
			return
		}

		idAgence, errAgence := middleware.GetIDAgenceUtilisateur(database, idUtilisateurAdmin)
		if errAgence != nil {
			fmt.Printf("Erreur récupération agence : %v\n", errAgence)
			http.Error(response, "Agence introuvable pour cet utilisateur", http.StatusForbidden)
			return
		}

		if errUpdate := MiseAJourPlanningsExpires(database); errUpdate != nil {
			fmt.Printf("Erreur mise à jour planning expirés : %v", errUpdate)
		}

		rowsPlanning, rowsErr := database.Query(`SELECT p.id_planning, u.nom, u.prenom, p.date, p.heure_debut, p.heure_fin,  p.statut FROM planning p
												JOIN benevole b ON b.id_benevole = p.id_benevole
												JOIN utilisateur u ON u.id_utilisateur = b.id_utilisateur
												WHERE u.id_agence = $1 ORDER BY p.date, p.heure_debut 
		`, idAgence)

		if rowsErr != nil {
			fmt.Printf("Voici l'erreur à la ligne 372 : %v\n", rowsErr)
			http.Error(response, "Erreur lors de la récupération des planning", http.StatusInternalServerError)
			return
		}
		defer rowsPlanning.Close()

		var planning_List []models.PlanningAfficheDashboard

		for rowsPlanning.Next() {
			var planning models.PlanningAfficheDashboard

			errQuery := rowsPlanning.Scan(&planning.ID_Planning,
				&planning.Nom,
				&planning.Prenom,
				&planning.Date,
				&planning.Heure_Debut,
				&planning.Heure_Fin,
				&planning.Statut,
			)

			if errQuery != nil {
				fmt.Printf("Erreur lors du Scan : %v\n", errQuery)
				http.Error(response, "Erreur lors de la lecture des données", http.StatusInternalServerError)
				return
			}

			planning_List = append(planning_List, planning)
		}

		data := AdminAgenceDashboardPlanningAfficher{
			Planning: planning_List,
		}

		tmpl, errTmpl := template.ParseFiles("./templates/admin_agence/planning.html")
		if errTmpl != nil {
			fmt.Printf("Erreur chargement template : %v\n", errTmpl)
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "ADMIN_AGENCE")
}

func FormModifierPlanning(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		idPlanning := request.URL.Query().Get("id")
		if idPlanning == "" {
			http.Error(response, "ID manquant", http.StatusBadRequest)
			return
		}

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur lors de la récupération de la session", http.StatusInternalServerError)
			return
		}

		idUtilisateurAdmin, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			http.Error(response, "Utilisateur non identifié", http.StatusUnauthorized)
			return
		}

		idAgence, agenceErr := middleware.GetIDAgenceUtilisateur(database, idUtilisateurAdmin)
		if agenceErr != nil {
			http.Error(response, "Agence introuvable pour cet utilisateur", http.StatusForbidden)
			return
		}

		var planning models.Planning
		errQuery := database.QueryRow(`SELECT p.id_planning, p.id_benevole, p.id_disponibilite, p.date, p.heure_debut, p.heure_fin, p.statut
									FROM planning p
									JOIN benevole b ON b.id_benevole = p.id_benevole
									JOIN utilisateur u ON u.id_utilisateur = b.id_utilisateur
									WHERE p.id_planning = $1 AND u.id_agence = $2
		`, idPlanning, idAgence).Scan(&planning.ID_Planning, &planning.ID_Benevole, &planning.ID_Disponibilite, &planning.Date, &planning.Heure_Debut, &planning.Heure_Fin, &planning.Statut)

		if errQuery == sql.ErrNoRows {
			http.Error(response, "Planning introuvable dans votre agence", http.StatusForbidden)
			return
		}

		if errQuery != nil {
			fmt.Printf("Erreur : %v\n", errQuery)
			http.Error(response, "Erreur select dans la base de données", http.StatusInternalServerError)
			return
		}

		data := AdminAgenceDashboardPlanning{
			Planning: planning,
		}

		tmpl, errTmpl := template.ParseFiles("./templates/admin_agence/modifier_planning.html")
		if errTmpl != nil {
			fmt.Printf("Erreur chargement template : %v\n", errTmpl)
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(response, data)

	}, "ADMIN_AGENCE")
}

func ModifierPlanning(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur lors de la récupération de la session", http.StatusInternalServerError)
			return
		}

		idUtilisateurAdmin, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			http.Error(response, "Utilisateur non identifié", http.StatusUnauthorized)
			return
		}

		idAgence, errAgence := middleware.GetIDAgenceUtilisateur(database, idUtilisateurAdmin)
		if errAgence != nil {
			http.Error(response, "Agence introuvable pour cet utilisateur", http.StatusForbidden)
			return
		}

		errForm := request.ParseForm()
		if errForm != nil {
			http.Error(response, "Données de formulaire invalides", http.StatusBadRequest)
			return
		}

		idPlanning := request.FormValue("id_planning")
		if idPlanning == "" {
			http.Error(response, "ID manquant", http.StatusBadRequest)
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

		result, errExec := database.Exec(`
			UPDATE planning p SET date = $1, heure_debut = $2, heure_fin = $3, statut = $4 FROM benevole b
			JOIN utilisateur u ON b.id_utilisateur = u.id_utilisateur
			WHERE p.id_planning = $5
				AND p.id_benevole = b.id_benevole
				AND u.id_agence = $6;
		`, date, heure_debut, heure_fin, statut, idPlanning, idAgence)

		if errExec != nil {
			fmt.Printf("Erreur : %v", errExec)
			http.Error(response, "Erreur lors de l'update du planning", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(response, "Planning introuvable dans votre agence", http.StatusForbidden)
			return
		}

		fmt.Println("Planning modifié avec succès !")
		http.Redirect(response, request, "/admin-agence/planning", http.StatusSeeOther)

	}, "ADMIN_AGENCE")
}

func DeletePlanning(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			fmt.Printf("Erreur récupération session : %v\n", sessErr)
			http.Error(response, "Erreur récupération de session", http.StatusInternalServerError)
			return
		}

		idUtilisateurAdmin, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			http.Error(response, "Utilisateur non identifié", http.StatusUnauthorized)
			return
		}

		idAgence, errAgence := middleware.GetIDAgenceUtilisateur(database, idUtilisateurAdmin)
		if errAgence != nil {
			http.Error(response, "Agence introuvable pour cet utilisateur", http.StatusForbidden)
			return
		}

		errForm := request.ParseForm()
		if errForm != nil {
			http.Error(response, "Données de formulaire invalides", http.StatusBadRequest)
			return
		}

		idPlanning := request.FormValue("id_planning")
		if idPlanning == "" {
			http.Error(response, "ID manquant", http.StatusBadRequest)
			return
		}

		tx, errTx := database.Begin()
		if errTx != nil {
			fmt.Printf("Erreur ouverture transaction : %v\n", errTx)
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		var idDisponibilite int
		errQuery := tx.QueryRow(`
			SELECT p.id_disponibilite
			FROM planning p
			JOIN benevole b ON b.id_benevole = p.id_benevole
			JOIN utilisateur u ON u.id_utilisateur = b.id_utilisateur
			WHERE p.id_planning = $1 AND u.id_agence = $2
		`, idPlanning, idAgence).Scan(&idDisponibilite)

		if errQuery == sql.ErrNoRows {
			http.Error(response, "Planning introuvable dans votre agence", http.StatusForbidden)
			return
		}
		if errQuery != nil {
			fmt.Printf("Erreur récupération disponibilité : %v\n", errQuery)
			http.Error(response, "Erreur récupération disponibilité", http.StatusInternalServerError)
			return
		}

		_, deleteErr := tx.Exec(`DELETE FROM planning WHERE id_planning = $1`, idPlanning)
		if deleteErr != nil {
			fmt.Printf("Erreur suppression planning : %v\n", deleteErr)
			http.Error(response, "Erreur suppression planning", http.StatusInternalServerError)
			return
		}

		_, changeErr := tx.Exec(`UPDATE disponibilite SET statut = 'DISPONIBLE' WHERE id_disponibilite = $1`, idDisponibilite)
		if changeErr != nil {
			fmt.Printf("Erreur mise à jour disponibilité : %v\n", changeErr)
			http.Error(response, "Erreur UPDATE de la disponibilité", http.StatusInternalServerError)
			return
		}

		if errCommit := tx.Commit(); errCommit != nil {
			fmt.Printf("Erreur commit transaction : %v\n", errCommit)
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}

		fmt.Println("Planning supprimé et disponibilité mise à jour")
		http.Redirect(response, request, "/admin-agence/planning", http.StatusSeeOther)

	}, "ADMIN_AGENCE", "ADMIN_GENERAL")
}

func DashboardAdminAgenceServices(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur récupération session", http.StatusInternalServerError)
			return
		}

		idUtilisateurAdmin, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			http.Error(response, "ID manquant", http.StatusInternalServerError)
			return
		}

		idAgence, errAgence := middleware.GetIDAgenceUtilisateur(database, idUtilisateurAdmin)
		if errAgence != nil {
			http.Error(response, "Agence introuvable pour cet utilisateur", http.StatusForbidden)
			return
		}

		rowsServices, servicesErr := database.Query(`SELECT s.id_service, s.nom, s.description, c.nom AS competence, s.statut FROM service s
													LEFT JOIN competence c ON s.id_competence = c.id_competence
													WHERE s.id_agence = $1
													ORDER BY s.id_service
		`, idAgence)
		if servicesErr != nil {
			fmt.Printf("Erreur : %v", servicesErr)
			http.Error(response, "Erreur récupération des services", http.StatusInternalServerError)
			return
		}
		defer rowsServices.Close()

		var services_List []models.Service

		for rowsServices.Next() {

			var service models.Service

			var competence sql.NullString

			errQuery := rowsServices.Scan(&service.ID_Service,
				&service.Nom,
				&service.Description,
				&competence,
				&service.Statut,
			)
			if errQuery != nil {
				fmt.Printf("Erreur : %v", errQuery)
				http.Error(response, "Erreur Scan service", http.StatusInternalServerError)
				return
			}

			if competence.Valid {
				service.Competence = competence.String
			} else {
				service.Competence = "Aucune compétence requise"
			}

			services_List = append(services_List, service)
		}

		rowsCompetences, competenceErr := database.Query(`SELECT id_competence, nom FROM competence`)
		if competenceErr != nil {
			fmt.Printf("Erreur : %v", competenceErr)
			http.Error(response, "Erreur récupération de compétences", http.StatusInternalServerError)
			return
		}
		defer rowsCompetences.Close()

		var competences_List []models.Competence

		for rowsCompetences.Next() {

			var competence models.Competence

			errQuery := rowsCompetences.Scan(&competence.IDCompetence, &competence.Nom)
			if errQuery != nil {
				fmt.Printf("Erreur : %v", errQuery)
				http.Error(response, "Erreur scan compétence", http.StatusInternalServerError)
				return
			}

			competences_List = append(competences_List, competence)
		}

		data := AdminAgenceDashboardServicesAfficher{
			Services:    services_List,
			Competences: competences_List,
		}

		tmpl, tmplErr := template.ParseFiles("./templates/admin_agence/services.html")
		if tmplErr != nil {
			fmt.Printf("Erreur : %v", tmplErr)
			http.Error(response, "Erreur parsage fichier", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "ADMIN_AGENCE")
}

func PageCreerService(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		rowsCompetences, errQuery := database.Query(`SELECT id_competence, nom FROM competence ORDER BY nom`)
		if errQuery != nil {
			fmt.Printf("Erreur : %v", errQuery)
			http.Error(response, "Erreur récupération des compétences", http.StatusInternalServerError)
			return
		}
		defer rowsCompetences.Close()

		var competencesList []models.Competence
		for rowsCompetences.Next() {
			var competence models.Competence
			if errScan := rowsCompetences.Scan(&competence.IDCompetence, &competence.Nom); errScan != nil {
				fmt.Printf("Erreur : %v", errScan)
				http.Error(response, "Erreur scan compétence", http.StatusInternalServerError)
				return
			}
			competencesList = append(competencesList, competence)
		}

		data := ServiceCreerData{Competences: competencesList}

		tmpl, errTmpl := template.ParseFiles("./templates/admin_agence/creer_service.html")
		if errTmpl != nil {
			fmt.Printf("Erreur : %v", errTmpl)
			http.Error(response, "Erreur parsage fichier", http.StatusInternalServerError)
			return
		}

		if errExec := tmpl.Execute(response, data); errExec != nil {
			fmt.Printf("Erreur exécution template : %v", errExec)
		}

	}, "ADMIN_AGENCE")
}

func CreateService(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur récupération de la session", http.StatusInternalServerError)
			return
		}

		idUtilisateurAdmin, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			http.Error(response, "ID manquant", http.StatusInternalServerError)
			return
		}

		idAgence, errAgence := middleware.GetIDAgenceUtilisateur(database, idUtilisateurAdmin)
		if errAgence != nil {
			http.Error(response, "Agence introuvable pour cet utilisateur", http.StatusForbidden)
			return
		}

		errForm := request.ParseForm()
		if errForm != nil {
			fmt.Printf("Erreur : %v", errForm)
			http.Error(response, "Erreur formulaire", http.StatusInternalServerError)
			return
		}

		nom := request.FormValue("nom")
		description := request.FormValue("description")
		competence := request.FormValue("competence")
		statut := request.FormValue("statut")

		var idCompetence sql.NullInt64

		if competence != "" {
			idConv, errConv := strconv.Atoi(competence)
			if errConv != nil {
				fmt.Printf("Erreur conversion idCompetence : %v", errConv)
				http.Error(response, "Compétence invalide", http.StatusBadRequest)
				return
			}
			idCompetence = sql.NullInt64{Int64: int64(idConv), Valid: true}
		}

		_, errInsert := database.Exec(`INSERT INTO service (nom, description, id_competence, statut, id_agence) VALUES ($1, $2, $3, $4, $5)`, nom, description, idCompetence, statut, idAgence)
		if errInsert != nil {
			fmt.Printf("Erreur : %v", errInsert)
			http.Error(response, "Erreur insert service", http.StatusInternalServerError)
			return
		}

		fmt.Println("Service crée avec succès !")
		http.Redirect(response, request, "/admin-agence/services", http.StatusSeeOther)

	}, "ADMIN_AGENCE")
}

func FormModifyService(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur récupération session", http.StatusInternalServerError)
			return
		}

		idUtilisateurAdmin, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			http.Error(response, "ID manquant", http.StatusInternalServerError)
			return
		}

		idAgence, errAgence := middleware.GetIDAgenceUtilisateur(database, idUtilisateurAdmin)
		if errAgence != nil {
			http.Error(response, "Agence introuvable pour cet utilisateur", http.StatusForbidden)
			return
		}

		idServiceStr := request.URL.Query().Get("id")
		if idServiceStr == "" {
			http.Error(response, "ID service manquant", http.StatusBadRequest)
			return
		}

		idService, errConv := strconv.Atoi(idServiceStr)
		if errConv != nil {
			http.Error(response, "ID service invalide", http.StatusBadRequest)
			return
		}

		var service models.ServiceModifierData
		var description sql.NullString
		var idCompetence sql.NullInt64

		errRow := database.QueryRow(`
			SELECT id_service, nom, description, id_competence, statut
			FROM service
			WHERE id_service = $1 AND id_agence = $2
		`, idService, idAgence).Scan(&service.IDService, &service.Nom, &description, &idCompetence, &service.Statut)

		if errRow == sql.ErrNoRows {
			http.Error(response, "Service introuvable", http.StatusNotFound)
			return
		}
		if errRow != nil {
			fmt.Printf("Erreur : %v", errRow)
			http.Error(response, "Erreur récupération service", http.StatusInternalServerError)
			return
		}

		if description.Valid {
			service.Description = description.String
		}
		if idCompetence.Valid {
			service.IDCompetence = int(idCompetence.Int64)
		}

		rowsCompetences, errQuery := database.Query(`SELECT id_competence, nom FROM competence ORDER BY nom`)
		if errQuery != nil {
			fmt.Printf("Erreur : %v", errQuery)
			http.Error(response, "Erreur récupération des compétences", http.StatusInternalServerError)
			return
		}
		defer rowsCompetences.Close()

		for rowsCompetences.Next() {

			var competence models.Competence

			if errScan := rowsCompetences.Scan(&competence.IDCompetence, &competence.Nom); errScan != nil {
				fmt.Printf("Erreur : %v", errScan)
				http.Error(response, "Erreur scan compétence", http.StatusInternalServerError)
				return
			}
			service.Competences = append(service.Competences, competence)
		}

		tmpl, errTmpl := template.ParseFiles("./templates/admin_agence/modifier_service.html")
		if errTmpl != nil {
			fmt.Printf("Erreur : %v", errTmpl)
			http.Error(response, "Erreur parsage fichier", http.StatusInternalServerError)
			return
		}

		if errExec := tmpl.Execute(response, service); errExec != nil {
			fmt.Printf("Erreur exécution template : %v", errExec)
		}

	}, "ADMIN_AGENCE")
}

func ModifierService(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur récupération de la session", http.StatusInternalServerError)
			return
		}

		idUtilisateurAdmin, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			http.Error(response, "ID manquant", http.StatusInternalServerError)
			return
		}

		idAgence, errAgence := middleware.GetIDAgenceUtilisateur(database, idUtilisateurAdmin)
		if errAgence != nil {
			fmt.Printf("Erreur : %v", errAgence)
			http.Error(response, "Agence introuvable pour cet utilisateur", http.StatusForbidden)
			return
		}

		errForm := request.ParseForm()
		if errForm != nil {
			fmt.Printf("Erreur : %v", errForm)
			http.Error(response, "Erreur formulaire", http.StatusInternalServerError)
			return
		}

		idService := request.FormValue("id_service")
		if idService == "" {
			http.Error(response, "ID service manquant", http.StatusBadRequest)
			return
		}

		nom := request.FormValue("nom")
		description := request.FormValue("description")
		competence := request.FormValue("competence")
		statut := request.FormValue("statut")

		var idCompetence sql.NullInt64
		if competence != "" {
			idConv, errConv := strconv.Atoi(competence)
			if errConv != nil {
				fmt.Printf("Erreur conversion idCompetence : %v", errConv)
				http.Error(response, "Compétence invalide", http.StatusBadRequest)
				return
			}
			idCompetence = sql.NullInt64{Int64: int64(idConv), Valid: true}
		}

		_, errExec := database.Exec(`UPDATE service SET 
										nom = COALESCE(NULLIF($1, ''), nom),
										description = COALESCE(NULLIF($2, ''), description),
										id_competence = $3,
										statut = COALESCE(NULLIF($4, ''), statut)
									WHERE id_service = $5 AND id_agence = $6
		`, nom, description, idCompetence, statut, idService, idAgence)
		if errExec != nil {
			fmt.Printf("Erreur : %v", errExec)
			http.Error(response, "Erreur modification service", http.StatusInternalServerError)
			return
		}

		http.Redirect(response, request, "/admin-agence/services", http.StatusSeeOther)

	}, "ADMIN_AGENCE")
}

func DeleteService(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			fmt.Printf("Erreur récupération session : %v\n", sessErr)
			http.Error(response, "Erreur récupération de session", http.StatusInternalServerError)
			return
		}

		idUtilisateurAdmin, ok := sess.Values["id_utilisateur"].(int)
		if !ok {
			http.Error(response, "Utilisateur non identifié", http.StatusUnauthorized)
			return
		}

		idAgence, errAgence := middleware.GetIDAgenceUtilisateur(database, idUtilisateurAdmin)
		if errAgence != nil {
			http.Error(response, "Agence introuvable pour cet utilisateur", http.StatusForbidden)
			return
		}

		errForm := request.ParseForm()
		if errForm != nil {
			http.Error(response, "Données de formulaire invalides", http.StatusBadRequest)
			return
		}

		idService := request.FormValue("id_service")
		if idService == "" {
			http.Error(response, "ID manquant", http.StatusBadRequest)
			return
		}

		fmt.Println("ID Service :", idService)
		fmt.Println("ID Agence :", idAgence)

		result, errExec := database.Exec(`DELETE FROM service WHERE id_service = $1 AND id_agence = $2`, idService, idAgence)
		if errExec != nil {
			fmt.Printf("Erreur : %v", errExec)
			http.Error(response, "Erreur suppression service", http.StatusInternalServerError)
			return
		}

		nbLignes, err := result.RowsAffected()
		if err != nil {
			fmt.Printf("Erreur RowsAffected : %v\n", err)
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}

		if nbLignes == 0 {
			http.Error(response, "Aucun service trouvé", http.StatusNotFound)
			return
		}

		fmt.Println("Service supprimé avec succès !")
		response.WriteHeader(http.StatusOK)

	}, "ADMIN_AGENCE")
}
