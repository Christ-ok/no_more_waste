package routes

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"no_more_waste/middleware"
	"no_more_waste/models"
	"no_more_waste/session"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
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

type AdminAgenceDashboardDemandeService struct {
	Demandes []models.DemandeService
}

type AdminAgenceBenevolesDisponibilites struct {
	Benevoles_Disponibilites []models.DemandeServiceDashboard
	ID_Demande               int
}

type AdminAgenceBenevoleDocument struct {
	Benevoles_Documents []models.BenevoleDocument
}

type AdminAgenceCollecteDashboard struct {
	Collecte []models.Collecte
}

type AdminAgenceBenevolesDisponibilitesCollect struct {
	Benevole_Disponibilites []models.DemandeCollecteDashboard
	ID_Collecte             int
}

type AdminAgenceStockDashboard struct {
	Stock []models.StockDashboard
}

type AdminAgenceTourneeData struct {
	Tournee []models.TourneeDashboard
}

type AdminAgenceCreerTournee struct {
	Stock []models.StockDisponible
}

type AdminAgenceBenevolesDisponibilitesTournee struct {
	Benevole_Disponibilites []models.TourneeDashboardAffectation
	ID_Tournee              int
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

func DashboardAdminAgenceDemandesServices(database *sql.DB) http.HandlerFunc {
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

		rowsDemandes, demandesErr := database.Query(`SELECT ds.id_demande_service, s.nom, u.nom, u.prenom, ds.date_demande, ds.statut FROM demande_service ds
														JOIN service s ON ds.id_service = s.id_service
														JOIN adherent a ON ds.id_adherent = a.id_adherent
														JOIN utilisateur u ON a.id_utilisateur = u.id_utilisateur 
													WHERE s.id_agence = $1
													ORDER BY ds.date_demande DESC
		`, idAgence)
		if demandesErr != nil {
			fmt.Printf("Erreur : %v", demandesErr)
			http.Error(response, "Erreur récupération demandes", http.StatusInternalServerError)
			return
		}
		defer rowsDemandes.Close()

		var demandes_List []models.DemandeService

		for rowsDemandes.Next() {

			var demande models.DemandeService

			errQuery := rowsDemandes.Scan(&demande.ID_Demande_Service,
				&demande.Nom_Service,
				&demande.Nom_Adherent,
				&demande.Prenom_Adherent,
				&demande.Date_Demande,
				&demande.Statut,
			)
			if errQuery != nil {
				fmt.Printf("Erreur : %v", errQuery)
				http.Error(response, "Erreur Scan demandes", http.StatusInternalServerError)
				return
			}

			demandes_List = append(demandes_List, demande)
		}

		data := AdminAgenceDashboardDemandeService{
			Demandes: demandes_List,
		}

		tmpl, tmplErr := template.ParseFiles("./templates/admin_agence/demandes_services.html")
		if tmplErr != nil {
			fmt.Printf("Erreur : %v", tmplErr)
			http.Error(response, "Erreur parse fichier demandes_services.html", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "ADMIN_AGENCE")
}

func PageAffectationBenevoleService(database *sql.DB) http.HandlerFunc {
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
			fmt.Printf("Erreur : %v", errForm)
			http.Error(response, "Erreur parse formulaire", http.StatusInternalServerError)
			return
		}

		idDemandeStr := request.FormValue("id_demande")

		idDemande, errDemande := strconv.Atoi(idDemandeStr)
		if errDemande != nil {
			fmt.Printf("Erreur : %v", errDemande)
			http.Error(response, "Erreur conversion id", http.StatusInternalServerError)
			return
		}

		var idCompetenceRequise int
		errCompetence := database.QueryRow(`SELECT s.id_competence FROM demande_service ds JOIN service s ON s.id_service = ds.id_service WHERE ds.id_demande_service = $1`, idDemande).Scan(&idCompetenceRequise)
		if errCompetence != nil {
			fmt.Printf("Erreur : %v", errCompetence)
			http.Error(response, "Erreur récupération ID compétence", http.StatusInternalServerError)
			return
		}

		rowsBenevoleDispo, errBenevole := database.Query(`SELECT b.id_benevole, u.nom, u.prenom, p.id_planning, p.date, p.heure_debut, p.heure_fin, p.statut FROM benevole b
														  	JOIN utilisateur u ON u.id_utilisateur = b.id_utilisateur
														  	JOIN planning p ON p.id_benevole = b.id_benevole
														  WHERE u.id_agence = $1
														  	AND b.id_competence = $2
															AND b.statut = 'VALIDE'
															AND p.statut = 'PLANIFIE'
														  ORDER BY p.date, p.heure_debut
		`, idAgence, idCompetenceRequise)
		if errBenevole != nil {
			fmt.Printf("Erreur : %v", errBenevole)
			http.Error(response, "Erreur récupération bénévoles disponibles", http.StatusInternalServerError)
			return
		}
		defer rowsBenevoleDispo.Close()

		var benevolesDisponibilites_List []models.DemandeServiceDashboard

		for rowsBenevoleDispo.Next() {

			var benevoleDisponibilite models.DemandeServiceDashboard

			errScan := rowsBenevoleDispo.Scan(&benevoleDisponibilite.ID_Benevole,
				&benevoleDisponibilite.Nom,
				&benevoleDisponibilite.Prenom,
				&benevoleDisponibilite.ID_Planning,
				&benevoleDisponibilite.Date_Planning,
				&benevoleDisponibilite.Heure_Debut,
				&benevoleDisponibilite.Heure_Fin,
				&benevoleDisponibilite.Statut,
			)
			if errScan != nil {
				fmt.Printf("Erreur : %v", errScan)
				http.Error(response, "Erreur scan bénévoles disponibilités", http.StatusInternalServerError)
				return
			}

			benevolesDisponibilites_List = append(benevolesDisponibilites_List, benevoleDisponibilite)
		}

		data := AdminAgenceBenevolesDisponibilites{
			Benevoles_Disponibilites: benevolesDisponibilites_List,
			ID_Demande:               idDemande,
		}

		tmpl, tmplErr := template.ParseFiles("./templates/admin_agence/benevoles_disponibilites.html")
		if tmplErr != nil {
			fmt.Printf("Erreur : %v", tmplErr)
			http.Error(response, "Erreur pase file html", http.StatusInternalServerError)
			return
		}

		if errExec := tmpl.Execute(response, data); errExec != nil {
			fmt.Printf("Erreur exécutio template : %v", errExec)
		}

	}, "ADMIN_AGENCE")
}

func DashboardAdminAgenceDocumentsBenevoles(database *sql.DB) http.HandlerFunc {
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

		rows, errRows := database.Query(`SELECT b.id_benevole, u.nom, u.prenom, c.nom , b.statut, b.justificatif FROM benevole b
											JOIN utilisateur u ON u.id_utilisateur = b.id_utilisateur
											LEFT JOIN competence c ON c.id_competence = b.id_competence
										WHERE u.id_agence = $1	
		`, idAgence)
		if errRows != nil {
			fmt.Printf("Erreur : %v", errRows)
			http.Error(response, "Erreur récupération documents", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var benevolesDocuments_List []models.BenevoleDocument

		for rows.Next() {
			var benevoleDocument models.BenevoleDocument
			var justificatif sql.NullString

			errQuery := rows.Scan(&benevoleDocument.ID_Benevole,
				&benevoleDocument.Nom,
				&benevoleDocument.Prenom,
				&benevoleDocument.Nom_Competence,
				&benevoleDocument.Statut,
				&justificatif,
			)
			if errQuery != nil {
				fmt.Printf("Erreur : %v", errQuery)
				http.Error(response, "Erreur Scan", http.StatusInternalServerError)
				return
			}

			if justificatif.Valid {
				benevoleDocument.Justificatif = justificatif.String
			}

			benevolesDocuments_List = append(benevolesDocuments_List, benevoleDocument)
		}

		data := AdminAgenceBenevoleDocument{
			Benevoles_Documents: benevolesDocuments_List,
		}

		tmpl, errTmpl := template.ParseFiles("./templates/admin_agence/benevoles_documents.html")
		if errTmpl != nil {
			fmt.Printf("Erreur : %v", errTmpl)
			http.Error(response, "Erreur parse file html", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "ADMIN_AGENCE")
}

func VoirDocumentBenevole(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		id := request.URL.Query().Get("id")
		if id == "" {
			http.Error(response, "ID manquant", http.StatusInternalServerError)
			return
		}

		var chemin string

		err := database.QueryRow(`
            SELECT justificatif
            FROM benevole
            WHERE id_benevole = $1
        `, id).Scan(&chemin)

		if err != nil {
			http.Error(response, "Document introuvable", http.StatusNotFound)
			return
		}

		http.ServeFile(response, request, chemin)

	}, "ADMIN_AGENCE")
}

func ValiderBenevole(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		if err := request.ParseForm(); err != nil {
			http.Error(response, "Requête invalide", http.StatusBadRequest)
			return
		}

		idBenevole := request.FormValue("id_benevole")
		if idBenevole == "" {
			http.Error(response, "ID bénévole manquant", http.StatusBadRequest)
			return
		}

		_, err := database.Exec(`
			UPDATE benevole
			SET statut = 'VALIDE'
			WHERE id_benevole = $1
		`, idBenevole)

		if err != nil {
			fmt.Printf("Erreur validation bénévole : %v\n", err)
			http.Error(response, "Erreur lors de la validation", http.StatusInternalServerError)
			return
		}

		http.Redirect(response, request, "/admin-agence/benevoles/documents", http.StatusSeeOther)

	}, "ADMIN_AGENCE")
}

func RejeterBenevole(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		if err := request.ParseForm(); err != nil {
			http.Error(response, "Requête invalide", http.StatusBadRequest)
			return
		}

		idBenevole := request.FormValue("id_benevole")
		if idBenevole == "" {
			http.Error(response, "ID bénévole manquant", http.StatusBadRequest)
			return
		}

		_, err := database.Exec(`
			UPDATE benevole
			SET statut = 'REFUSE'
			WHERE id_benevole = $1
		`, idBenevole)

		if err != nil {
			fmt.Printf("Erreur rejet bénévole : %v\n", err)
			http.Error(response, "Erreur lors du rejet", http.StatusInternalServerError)
			return
		}

		http.Redirect(response, request, "/admin-agence/benevoles/documents", http.StatusSeeOther)

	}, "ADMIN_AGENCE")
}

func AttributionDemandeServicePlanningBenevole(database *sql.DB) http.HandlerFunc {
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
			fmt.Printf("Erreur : %v", errForm)
			http.Error(response, "Erreur formulaire", http.StatusInternalServerError)
			return
		}

		fmt.Println("id_demande :", request.FormValue("id_demande"))
		fmt.Println("id_planning :", request.FormValue("id_planning"))
		fmt.Println("id_benevole :", request.FormValue("id_benevole"))

		idPlanning, errPlanning := strconv.Atoi(request.FormValue("id_planning"))
		if errPlanning != nil {
			fmt.Printf("Erreur : %v", errPlanning)
			http.Error(response, "ID Planning vide", http.StatusInternalServerError)
			return
		}

		idBenevole, errBenevole := strconv.Atoi(request.FormValue("id_benevole"))
		if errBenevole != nil {
			fmt.Printf("Erreur : %v", errBenevole)
			http.Error(response, "ID Benevole vide", http.StatusInternalServerError)
			return
		}

		idDemande, errDemande := strconv.Atoi(request.FormValue("id_demande"))
		if errDemande != nil {
			fmt.Printf("Erreur : %v", errDemande)
			http.Error(response, "ID demande vide", http.StatusBadRequest)
			return
		}

		tx, txErr := database.Begin()
		if txErr != nil {
			fmt.Printf("Erreur : %v", txErr)
			http.Error(response, "Erreur ouverture transaction", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		result, errExec := tx.Exec(`UPDATE demande_service ds SET 
										id_benevole = $1,
										id_planning = $2,
										statut = 'ATTRIBUE'
									 FROM service s
										WHERE ds.id_service = s.id_service
										AND ds.id_demande_service = $3
										AND s.id_agence = $4
		`, idBenevole, idPlanning, idDemande, idAgence)
		if errExec != nil {
			fmt.Printf("Erreur : %v", errExec)
			http.Error(response, "Erreur attribution planning et bénévole", http.StatusInternalServerError)
			return
		}

		rows, _ := result.RowsAffected()
		if rows == 0 {
			http.Error(response, "Demande introuvable", http.StatusNotFound)
			return
		}

		resultPlanning, errUpdatePlanning := tx.Exec(`UPDATE planning SET statut = 'ATTRIBUE' WHERE id_planning = $1`, idPlanning)
		if errUpdatePlanning != nil {
			fmt.Printf("Erreur : %v", errUpdatePlanning)
			http.Error(response, "Erreur update planning", http.StatusInternalServerError)
			return
		}

		rowsPlanning, _ := resultPlanning.RowsAffected()

		if rowsPlanning == 0 {
			http.Error(response, "Planning introuvable", http.StatusNotFound)
			return
		}

		if errCommit := tx.Commit(); errCommit != nil {
			fmt.Printf("Erreur : %v", errCommit)
			http.Error(response, "Erreur commit transaction", http.StatusInternalServerError)
			return
		}

		var nom_benevole, prenom_benevole, statut_service string
		var date_service, heure_debut, heure_fin time.Time

		errInfo := database.QueryRow(`SELECT u.nom, u.prenom, p.date, p.heure_debut, p.heure_fin, p.statut FROM benevole b
										JOIN utilisateur u ON u.id_utilisateur = b.id_utilisateur
										JOIN planning p ON p.id_planning = $2
									  WHERE b.id_benevole = $1
		`, idBenevole, idPlanning).Scan(&nom_benevole, &prenom_benevole, &date_service, &heure_debut, &heure_fin, &statut_service)
		if errInfo != nil {
			fmt.Printf("Erreur : %v", errInfo)
			http.Redirect(response, request, "/admin-agence/demande-services", http.StatusSeeOther)
			return
		}

		if err := genererPlanningExcel(database, idPlanning, idDemande, idBenevole, nom_benevole, prenom_benevole, date_service, heure_debut, heure_fin, statut_service); err != nil {
			fmt.Printf("Erreur génération planning excel : %v\n", err)
			return
		}

		fmt.Println("Demande de service attribué avec succès !")
		http.Redirect(response, request, "/admin-agence/demande-services", http.StatusSeeOther)

	}, "ADMIN_AGENCE")
}

func HistoriqueServicesRealises(database *sql.DB) http.HandlerFunc {
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

		rowsServices, errRows := database.Query(`SELECT ds.id_demande_service, s.nom, p.date, ua.nom, ua.prenom, ub.nom, ub.prenom, ds.statut FROM demande_service ds
												  JOIN service s ON s.id_service = ds.id_service
												  JOIN planning p ON p.id_planning = ds.id_planning
												  JOIN adherent a ON ds.id_adherent = a.id_adherent
												  JOIN utilisateur ua ON a.id_utilisateur = ua.id_utilisateur
												  LEFT JOIN benevole b ON ds.id_benevole = b.id_benevole
												  LEFT JOIN utilisateur ub ON b.id_utilisateur = ub.id_utilisateur
												WHERE ua.id_agence = $1
		`, idAgence)
		if errRows != nil {
			fmt.Printf("Erreur : %v", errRows)
			http.Error(response, "Erreur récupération Services", http.StatusInternalServerError)
			return
		}
		defer rowsServices.Close()

		var services_List []models.DemandeService

		for rowsServices.Next() {

			var service models.DemandeService

			errScan := rowsServices.Scan(&service.ID_Demande_Service,
				&service.Nom_Service,
				&service.Date_Demande,
				&service.Nom_Adherent,
				&service.Prenom_Adherent,
				&service.Nom_Benevole,
				&service.Prenom_Benevole,
				&service.Statut,
			)
			if errScan != nil {
				fmt.Printf("Erreur : %v", errScan)
				http.Error(response, "Erreur scan service", http.StatusInternalServerError)
				return
			}

			services_List = append(services_List, service)
		}

		data := AdminAgenceDashboardDemandeService{
			Demandes: services_List,
		}

		tmpl, errTmpl := template.ParseFiles("./templates/admin_agence/historiques_services.html")
		if errTmpl != nil {
			fmt.Printf("Erreur : %v", errTmpl)
			http.Error(response, "Erreur parseFile html", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "ADMIN_AGENCE")
}

func genererPlanningExcel(database *sql.DB, idPlanning int, id_Demande int, idBenevole int,
	nom_benevole string, prenom_benevole string, date_service time.Time,
	heure_debut time.Time, heure_fin time.Time, statut_service string) error {

	var nom_service string
	err := database.QueryRow(`
		SELECT s.nom FROM demande_service ds
		JOIN service s ON ds.id_service = s.id_service
		WHERE ds.id_demande_service = $1
	`, id_Demande).Scan(&nom_service)
	if err != nil {
		return fmt.Errorf("récupération nom service : %w", err)
	}

	f := excelize.NewFile()
	sheet := "Planning"
	f.SetSheetName("Sheet1", sheet)

	f.SetCellValue(sheet, "A1", "Planning du service")
	f.MergeCell(sheet, "A1", "G1")

	styleTitre, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  14,
			Color: "33553A",
		},
	})

	f.SetCellStyle(sheet, "A1", "G1", styleTitre)

	f.SetCellValue(sheet, "A2", nom_service)
	f.MergeCell(sheet, "A2", "G2")

	headers := []string{"Service", "Nom", "Prénom", "Date", "Début", "Fin", "Statut"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 4)
		f.SetCellValue(sheet, cell, header)
	}

	styleHeader, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"33553A"}},
	})

	f.SetCellStyle(sheet, "A4", "G4", styleHeader)

	f.SetCellValue(sheet, "A5", nom_service)
	f.SetCellValue(sheet, "B5", nom_benevole)
	f.SetCellValue(sheet, "C5", prenom_benevole)
	f.SetCellValue(sheet, "D5", date_service.Format("02/01/2006"))
	f.SetCellValue(sheet, "E5", heure_debut.Format("15:04"))
	f.SetCellValue(sheet, "F5", heure_fin.Format("15:04"))
	f.SetCellValue(sheet, "G5", statut_service)

	f.SetColWidth(sheet, "A", "A", 25)
	f.SetColWidth(sheet, "B", "B", 18)
	f.SetColWidth(sheet, "C", "C", 18)
	f.SetColWidth(sheet, "D", "D", 14)
	f.SetColWidth(sheet, "E", "E", 12)
	f.SetColWidth(sheet, "F", "F", 12)
	f.SetColWidth(sheet, "G", "G", 15)

	cheminDossier := filepath.Join("stockage", "plannings")
	if err := os.MkdirAll(cheminDossier, 0755); err != nil {
		return fmt.Errorf("Création dossier stockage : %w", err)
	}

	nomFichier := fmt.Sprintf("planning_%d.xlsx", idPlanning)
	cheminFichier := filepath.Join(cheminDossier, nomFichier)

	if err := f.SaveAs(cheminFichier); err != nil {
		return fmt.Errorf("sauvegarde fichier excel : %w", err)
	}

	_, err = database.Exec(`
		INSERT INTO planning_excel (id_planning, id_benevole, chemin_fichier)
		VALUES ($1, $2, $3)
	`, idPlanning, idBenevole, cheminFichier)
	if err != nil {
		return fmt.Errorf("insertion planning_excel : %w", err)
	}

	return nil
}

func DashboardAdminAgenceCollecte(database *sql.DB) http.HandlerFunc {
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

		rowsCollecte, errCollecte := database.Query(`SELECT c.id_collecte, u.nom, u.prenom, c.date_collecte, c.statut FROM collecte c
													  JOIN commercant co ON c.id_commercant = co.id_commercant
													  JOIN utilisateur u ON u.id_utilisateur = co.id_utilisateur 
													 WHERE c.id_agence = $1
		`, idAgence)
		if errCollecte != nil {
			fmt.Printf("Erreur : %v", errCollecte)
			http.Error(response, "Erreur récupération collecte", http.StatusInternalServerError)
			return
		}
		defer rowsCollecte.Close()

		var collectes_List []models.Collecte

		for rowsCollecte.Next() {
			var collecte models.Collecte

			errScan := rowsCollecte.Scan(&collecte.ID_Collecte,
				&collecte.Nom_Commercant,
				&collecte.Prenom_Commercant,
				&collecte.Date_Collecte,
				&collecte.Statut,
			)
			if errScan != nil {
				fmt.Printf("Erreur : %v", errScan)
				http.Error(response, "Erreur scan collecte", http.StatusInternalServerError)
				return
			}

			collectes_List = append(collectes_List, collecte)
		}

		data := AdminAgenceCollecteDashboard{
			Collecte: collectes_List,
		}

		tmpl, errTmpl := template.ParseFiles("./templates/admin_agence/collectes.html")
		if errTmpl != nil {
			fmt.Printf("Erreur : %v", errTmpl)
			http.Error(response, "Erreur parsefiles html", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "ADMIN_AGENCE")
}

func PageAffectationBenevoleCollecte(database *sql.DB) http.HandlerFunc {
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

		if errForm := request.ParseForm(); errForm != nil {
			fmt.Printf("Erreur : %v", errForm)
			http.Error(response, "Erreur parse formulaire", http.StatusInternalServerError)
			return
		}

		idCollecteStr := request.FormValue("id_collecte")
		if idCollecteStr == "" {
			http.Error(response, "ID collecte manquant", http.StatusBadRequest)
			return
		}

		idCollecte, errConv := strconv.Atoi(idCollecteStr)
		if errConv != nil {
			fmt.Printf("Erreur : %v", errConv)
			http.Error(response, "Erreur conversion id_collecte", http.StatusInternalServerError)
			return
		}

		rowsBenevoleDispo, errBenevole := database.Query(`SELECT b.id_benevole, u.nom, u.prenom, p.id_planning, p.date, p.heure_debut, p.heure_fin, p.statut FROM benevole b
														  	JOIN utilisateur u ON u.id_utilisateur = b.id_utilisateur
														  	JOIN planning p ON p.id_benevole = b.id_benevole
															JOIN competence c ON b.id_competence = c.id_competence
														  WHERE u.id_agence = $1
														  	AND c.nom = 'Chauffeur'
															AND b.statut = 'VALIDE'
															AND p.statut = 'PLANIFIE'
														  ORDER BY p.date, p.heure_debut	
		`, idAgence)
		if errBenevole != nil {
			fmt.Printf("Erreur : %v", errBenevole)
			http.Error(response, "Erreur récupération bénévole", http.StatusInternalServerError)
			return
		}
		defer rowsBenevoleDispo.Close()

		var benevolesDisponibilites_List []models.DemandeCollecteDashboard

		for rowsBenevoleDispo.Next() {
			var benevoleDispo models.DemandeCollecteDashboard

			errScan := rowsBenevoleDispo.Scan(&benevoleDispo.ID_Benevole,
				&benevoleDispo.Nom_Benevole,
				&benevoleDispo.Prenom_Benevole,
				&benevoleDispo.ID_Planning,
				&benevoleDispo.Date_Planning,
				&benevoleDispo.Heure_Debut,
				&benevoleDispo.Heure_Fin,
				&benevoleDispo.Statut_Planning,
			)
			if errScan != nil {
				fmt.Printf("Erreur : %v", errScan)
				http.Error(response, "Erreur scan dispo benevole", http.StatusInternalServerError)
				return
			}

			benevolesDisponibilites_List = append(benevolesDisponibilites_List, benevoleDispo)
		}

		data := AdminAgenceBenevolesDisponibilitesCollect{
			Benevole_Disponibilites: benevolesDisponibilites_List,
			ID_Collecte:             idCollecte,
		}

		tmpl, tmplErr := template.ParseFiles("./templates/admin_agence/benevoles_disponibilites_collectes.html")
		if tmplErr != nil {
			fmt.Printf("Erreur : %v", tmplErr)
			http.Error(response, "Erreur parsefiles html", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "ADMIN_AGENCE")
}

func AttributionCollectePlanningBenevole(database *sql.DB) http.HandlerFunc {
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
			fmt.Printf("Erreur : %v", errForm)
			http.Error(response, "Erreur formulaire", http.StatusInternalServerError)
			return
		}

		idCollecte, errCollecte := strconv.Atoi(request.FormValue("id_collecte"))
		if errCollecte != nil {
			fmt.Printf("Erreur : %v", errCollecte)
			http.Error(response, "Erreur conversion collecte", http.StatusInternalServerError)
			return
		}

		idPlanning, errPlanning := strconv.Atoi(request.FormValue("id_planning"))
		if errPlanning != nil {
			fmt.Printf("Erreur : %v", errPlanning)
			http.Error(response, "ID Planning vide", http.StatusInternalServerError)
			return
		}

		idBenevole, errBenevole := strconv.Atoi(request.FormValue("id_benevole"))
		if errBenevole != nil {
			fmt.Printf("Erreur : %v", errBenevole)
			http.Error(response, "ID Benevole vide", http.StatusInternalServerError)
			return
		}

		tx, errTx := database.Begin()
		if errTx != nil {
			fmt.Printf("Erreur : %v", errTx)
			http.Error(response, "Erreur ouverture transaction", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		result, errExec := tx.Exec(`UPDATE collecte SET
										id_benevole = $1,
										id_planning = $2,
										statut = 'planifiee'
									WHERE id_collecte = $3
									AND id_agence = $4 
		`, idBenevole, idPlanning, idCollecte, idAgence)
		if errExec != nil {
			fmt.Printf("Erreur : %v", errExec)
			http.Error(response, "Erreur UPDATE table collecte", http.StatusInternalServerError)
			return
		}

		rows, _ := result.RowsAffected()
		if rows == 0 {
			http.Error(response, "Collecte introuvable", http.StatusNotFound)
			return
		}

		resultPlanning, errUpdatePlanning := tx.Exec(`UPDATE planning SET statut = 'ATTRIBUE' WHERE id_planning = $1`, idPlanning)
		if errUpdatePlanning != nil {
			fmt.Printf("Erreur : %v", errUpdatePlanning)
			http.Error(response, "Erreur Update statut Planning", http.StatusInternalServerError)
			return
		}

		rowsPlanning, _ := resultPlanning.RowsAffected()

		if rowsPlanning == 0 {
			http.Error(response, "Planning introuvable", http.StatusNotFound)
			return
		}

		if errCommit := tx.Commit(); errCommit != nil {
			fmt.Printf("Erreur : %v", errCommit)
			http.Error(response, "Erreur commit transaction", http.StatusInternalServerError)
			return
		}

		var nomBenevole, prenomBenevole, statutPlanning string
		var dateCollecte, heureDebut, heureFin time.Time

		errInfoPlanningExcel := database.QueryRow(`SELECT u.nom, u.prenom, p.date, p.heure_debut, p.heure_fin, p.statut FROM benevole b
													JOIN utilisateur u ON u.id_utilisateur = b.id_utilisateur
													JOIN planning p ON p.id_planning = $2
												   WHERE b.id_benevole = $1
		`, idBenevole, idPlanning).Scan(&nomBenevole, &prenomBenevole, &dateCollecte, &heureDebut, &heureFin, &statutPlanning)
		if errInfoPlanningExcel != nil {
			fmt.Printf("Erreur : %v", errInfoPlanningExcel)
			http.Error(response, "Erreur récupération infos planning", http.StatusInternalServerError)
			return
		}

		if err := genererPlanningExcelCollecte(database, idPlanning, idCollecte, idBenevole, nomBenevole, prenomBenevole, dateCollecte, heureDebut, heureFin, statutPlanning); err != nil {
			fmt.Printf("Erreur : %v", err)
			http.Error(response, "Collecte attribuée mais erreur lors de la génération du planning Excel", http.StatusInternalServerError)
			return
		}

		fmt.Println("Collecte attribuée avec succès !")
		http.Redirect(response, request, "/admin-agence/collectes", http.StatusSeeOther)

	}, "ADMIN_AGENCE")
}

func genererPlanningExcelCollecte(database *sql.DB, idPlanning int, idCollecte int, idBenevole int,
	nomBenevole string, prenomBenevole string, dateCollecte time.Time,
	heureDebut time.Time, heureFin time.Time, statutPlanning string) error {

	var nomCommercant, prenomCommercant, adresseCommercant string
	errCommercant := database.QueryRow(`SELECT u.nom, u.prenom, u.adresse FROM collecte c
											JOIN commercant co ON c.id_commercant = co.id_commercant
											JOIN utilisateur u ON u.id_utilisateur = co.id_utilisateur
										WHERE c.id_collecte = $1
	`, idCollecte).Scan(&nomCommercant, &prenomCommercant, &adresseCommercant)
	if errCommercant != nil {
		return fmt.Errorf("récupération infos commerçant : %w", errCommercant)
	}

	f := excelize.NewFile()
	sheet := "Planning"
	f.SetSheetName("Sheet1", sheet)

	f.SetCellValue(sheet, "A1", "Planning de collecte")
	f.MergeCell(sheet, "A1", "G1")

	styleTitre, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  14,
			Color: "33553A",
		},
	})
	f.SetCellStyle(sheet, "A1", "G1", styleTitre)

	f.SetCellValue(sheet, "A2", fmt.Sprintf("%s %s — %s", nomCommercant, prenomCommercant, adresseCommercant))
	f.MergeCell(sheet, "A2", "G2")

	headers := []string{"Commerçant", "Bénévole", "Prénom", "Date", "Début", "Fin", "Statut"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 4)
		f.SetCellValue(sheet, cell, header)
	}

	styleHeader, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"33553A"}},
	})
	f.SetCellStyle(sheet, "A4", "G4", styleHeader)

	f.SetCellValue(sheet, "A5", nomCommercant)
	f.SetCellValue(sheet, "B5", nomBenevole)
	f.SetCellValue(sheet, "C5", prenomBenevole)
	f.SetCellValue(sheet, "D5", dateCollecte.Format("02/01/2006"))
	f.SetCellValue(sheet, "E5", heureDebut.Format("15:04"))
	f.SetCellValue(sheet, "F5", heureFin.Format("15:04"))
	f.SetCellValue(sheet, "G5", statutPlanning)

	f.SetColWidth(sheet, "A", "A", 25)
	f.SetColWidth(sheet, "B", "B", 18)
	f.SetColWidth(sheet, "C", "C", 18)
	f.SetColWidth(sheet, "D", "D", 14)
	f.SetColWidth(sheet, "E", "E", 12)
	f.SetColWidth(sheet, "F", "F", 12)
	f.SetColWidth(sheet, "G", "G", 15)

	f.SetCellValue(sheet, "A7", "Produits à collecter")
	f.MergeCell(sheet, "A7", "G7")
	f.SetCellStyle(sheet, "A7", "G7", styleTitre)

	produitHeaders := []string{"Produit", "Quantité"}
	for i, header := range produitHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 8)
		f.SetCellValue(sheet, cell, header)
	}
	f.SetCellStyle(sheet, "A8", "B8", styleHeader)

	rowsProduits, errProduits := database.Query(`SELECT libelle, quantite FROM produit_collecte WHERE id_collecte = $1`, idCollecte)
	if errProduits != nil {
		return fmt.Errorf("récupération produits collecte : %w", errProduits)
	}
	defer rowsProduits.Close()

	ligne := 9
	for rowsProduits.Next() {
		var libelle string
		var quantite float64

		if errScan := rowsProduits.Scan(&libelle, &quantite); errScan != nil {
			return fmt.Errorf("scan produit collecte : %w", errScan)
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", ligne), libelle)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", ligne), quantite)
		ligne++
	}

	cheminDossier := filepath.Join("stockage", "plannings")
	if err := os.MkdirAll(cheminDossier, 0755); err != nil {
		return fmt.Errorf("création dossier stockage : %w", err)
	}

	nomFichier := fmt.Sprintf("planning_collecte_%d.xlsx", idPlanning)
	cheminFichier := filepath.Join(cheminDossier, nomFichier)

	if err := f.SaveAs(cheminFichier); err != nil {
		return fmt.Errorf("sauvegarde fichier excel : %w", err)
	}

	_, err := database.Exec(`
		INSERT INTO planning_excel (id_planning, id_benevole, chemin_fichier)
		VALUES ($1, $2, $3)
	`, idPlanning, idBenevole, cheminFichier)
	if err != nil {
		return fmt.Errorf("insertion planning_excel : %w", err)
	}

	return nil
}

func DashboardAdminAgenceStocks(database *sql.DB) http.HandlerFunc {
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

		rowsStocks, errStock := database.Query(`SELECT s.id_stock, pc.libelle, pc.code_barre, s.quantite_disponible, s.date_entree FROM stock s
													JOIN produit_collecte pc ON pc.id_stock = s.id_stock
												WHERE s.id_agence = $1
												ORDER BY s.date_entree DESC
		`, idAgence)
		if errStock != nil {
			fmt.Printf("Erreur : %v", errStock)
			http.Error(response, "Erreur récupération produits", http.StatusInternalServerError)
			return
		}
		defer rowsStocks.Close()

		var stock_List []models.StockDashboard

		for rowsStocks.Next() {
			var stock models.StockDashboard

			errScan := rowsStocks.Scan(&stock.ID_Stock,
				&stock.Libelle,
				&stock.Code_Barre,
				&stock.Quantite_Disponible,
				&stock.Date_Entree,
			)
			if errScan != nil {
				fmt.Printf("Erreur : %v", errScan)
				http.Error(response, "Erreur scan valeurs stock", http.StatusInternalServerError)
				return
			}

			stock_List = append(stock_List, stock)
		}

		data := AdminAgenceStockDashboard{
			Stock: stock_List,
		}

		tmpl, errTmpl := template.ParseFiles("./templates/admin_agence/stocks.html")
		if errTmpl != nil {
			fmt.Printf("Erreur : %v", errTmpl)
			http.Error(response, "Erreur parsefiles html", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "ADMIN_AGENCE")
}

func DashboardAdminAgenceTournee(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur récupération session", http.StatusInternalServerError)
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

		rowsTournee, errTournee := database.Query(`SELECT t.id_tournee, d.nom, d.type, t.date_tournee, t.statut, u.nom, u.prenom FROM tournee t
													JOIN destinataire d ON d.id_destinataire = t.id_destinataire
													LEFT JOIN benevole b ON b.id_benevole = t.id_benevole
													LEFT JOIN utilisateur u ON u.id_utilisateur = b.id_utilisateur
												   WHERE t.id_agence = $1
												   ORDER BY t.date_tournee
		`, idAgence)
		if errTournee != nil {
			fmt.Printf("Erreur : %v", errTournee)
			http.Error(response, "Erreur récupération tournee", http.StatusInternalServerError)
			return
		}
		defer rowsTournee.Close()

		var tournee_List []models.TourneeDashboard

		for rowsTournee.Next() {
			var tournee models.TourneeDashboard

			errScan := rowsTournee.Scan(&tournee.ID_Tournee,
				&tournee.Nom_Destinataire,
				&tournee.Type_Destinataire,
				&tournee.Date_Tournee,
				&tournee.Statut,
				&tournee.Nom_Benevole,
				&tournee.Prenom_Benevole,
			)
			if errScan != nil {
				fmt.Printf("Erreur : %v", errScan)
				http.Error(response, "Erreur Scan tournée", http.StatusInternalServerError)
				return
			}

			tournee_List = append(tournee_List, tournee)
		}

		data := AdminAgenceTourneeData{
			Tournee: tournee_List,
		}

		tmpl, errTmpl := template.ParseFiles("./templates/admin_agence/tournees.html")
		if errTmpl != nil {
			fmt.Printf("Erreur : %v", errTmpl)
			http.Error(response, "Erreur parsefile html", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "ADMIN_AGENCE")
}

func PageCreerTournee(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur récupération session", http.StatusInternalServerError)
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

		rowsStock, errStock := database.Query(`SELECT s.id_stock, pc.libelle, pc.code_barre, s.quantite_disponible
												FROM stock s
												JOIN produit_collecte pc ON pc.id_stock = s.id_stock
												WHERE s.id_agence = $1 AND s.quantite_disponible > 0
												ORDER BY pc.libelle
		`, idAgence)
		if errStock != nil {
			fmt.Printf("Erreur : %v", errStock)
			http.Error(response, "Erreur récupération stock", http.StatusInternalServerError)
			return
		}
		defer rowsStock.Close()

		var stock_List []models.StockDisponible

		for rowsStock.Next() {
			var stock models.StockDisponible

			errScan := rowsStock.Scan(&stock.ID_Stock, &stock.Libelle, &stock.Code_Barre, &stock.Quantite_Disponible)
			if errScan != nil {
				fmt.Printf("Erreur : %v", errScan)
				http.Error(response, "Erreur scan stock", http.StatusInternalServerError)
				return
			}

			stock_List = append(stock_List, stock)
		}

		data := AdminAgenceCreerTournee{
			Stock: stock_List,
		}

		tmpl, errTmpl := template.ParseFiles("./templates/admin_agence/creer_tournee.html")
		if errTmpl != nil {
			fmt.Printf("Erreur : %v", errTmpl)
			http.Error(response, "Erreur parsefiles html", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "ADMIN_AGENCE")
}

func CreerTournee(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur récupération session", http.StatusInternalServerError)
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
			fmt.Printf("Erreur : %v", errForm)
			http.Error(response, "Erreur formulaire", http.StatusInternalServerError)
			return
		}

		nomDestinataire := request.FormValue("nom_destinataire")
		adresseDestinataire := request.FormValue("adresse_destinataire")
		typeDestinataire := request.FormValue("type_destinataire")

		if nomDestinataire == "" || adresseDestinataire == "" || typeDestinataire == "" {
			http.Error(response, "Informations destinataire incomplètes", http.StatusBadRequest)
			return
		}

		dateTournee, errDate := time.Parse("2006-01-02", request.FormValue("date_tournee"))
		if errDate != nil {
			fmt.Printf("Erreur : %v", errDate)
			http.Error(response, "Format de date invalide", http.StatusBadRequest)
			return
		}

		idStock := request.Form["id_stock[]"]
		quantites := request.Form["quantite[]"]

		if len(idStock) == 0 || len(idStock) != len(quantites) {
			http.Error(response, "Liste de produits invalide", http.StatusBadRequest)
			return
		}

		tx, errTx := database.Begin()
		if errTx != nil {
			fmt.Printf("Erreur : %v", errTx)
			http.Error(response, "Erreur ouerture transaction", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		var idDestinataire int
		errInsertDestinataire := tx.QueryRow(`INSERT INTO destinataire (id_agence, type, nom, adresse)
											    VALUES ($1, $2, $3, $4) RETURNING id_destinataire
		`, idAgence, typeDestinataire, nomDestinataire, adresseDestinataire).Scan(&idDestinataire)
		if errInsertDestinataire != nil {
			fmt.Printf("Erreur : %v", errInsertDestinataire)
			http.Error(response, "Erreur Insertion Destinataire", http.StatusInternalServerError)
			return
		}

		var idTournee int
		errInsertTournee := tx.QueryRow(`INSERT INTO tournee (id_destinataire, id_agence, date_tournee, statut)
										  VALUES ($1, $2, $3, 'en_attente') RETURNING id_tournee
		`, idDestinataire, idAgence, dateTournee).Scan(&idTournee)
		if errInsertTournee != nil {
			fmt.Printf("Erreur : %v", errInsertTournee)
			http.Error(response, "Erreur Insertion Tournee", http.StatusInternalServerError)
			return
		}

		for i, idStockStr := range idStock {
			if idStockStr == "" {
				continue
			}

			idStock, errConvStock := strconv.Atoi(idStockStr)
			if errConvStock != nil {
				fmt.Printf("Erreur : %v", errConvStock)
				http.Error(response, "Erreur conversion idStock", http.StatusInternalServerError)
				return
			}

			quantite, errConvQuantite := strconv.ParseFloat(quantites[i], 64)
			if errConvQuantite != nil {
				fmt.Printf("Erreur : %v", errConvQuantite)
				http.Error(response, "Quantité invalide", http.StatusBadRequest)
				return
			}

			var quantiteDisponible float64
			errCheckStock := tx.QueryRow(`SELECT quantite_disponible FROM stock WHERE id_stock = $1 AND id_agence = $2`, idStock, idAgence).Scan(&quantiteDisponible)
			if errCheckStock != nil {
				fmt.Printf("Erreur : %v", errCheckStock)
				http.Error(response, "Produit en stock introuvable", http.StatusInternalServerError)
				return
			}

			if quantite > quantiteDisponible {
				http.Error(response, fmt.Sprintf("Quantité demandée (%.2f) supérieure au stock disponible (%.2f)", quantite, quantiteDisponible), http.StatusBadRequest)
				return
			}

			_, errInsertProduit := tx.Exec(`INSERT INTO produit_tournee (id_tournee, id_stock, quantite) VALUES ($1, $2, $3)`, idTournee, idStock, quantite)
			if errInsertProduit != nil {
				fmt.Printf("Erreur : %v", errInsertProduit)
				return
			}
		}

		if errCommit := tx.Commit(); errCommit != nil {
			fmt.Printf("Erreur : %v", errCommit)
			http.Error(response, "Erreur validation tournée", http.StatusInternalServerError)
			return
		}

		fmt.Println("Tournée créée avec succès !")
		http.Redirect(response, request, "/admin-agence/tournees", http.StatusSeeOther)

	}, "ADMIN_AGENCE")
}

func AfficherBenevoleDisponibleTournee(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur récupération session", http.StatusInternalServerError)
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

		rowsBenevole, errBenevole := database.Query(`SELECT b.id_benevole, u.nom, u.prenom, p.id_planning, p.date, p.heure_debut, p.heure_fin, p.statut FROM benevole b
														JOIN utilisateur u ON u.id_utilisateur = b.id_utilisateur
														JOIN planning p ON p.id_benevole = b.id_benevole
														JOIN competence c ON b.id_competence = c.id_competence
													 WHERE u.id_agence = $1
														AND c.nom = 'Chauffeur'
														AND b.statut = 'VALIDE'
														AND p.statut = 'PLANIFIE'
													 ORDER BY p.date, p.heure_debut								
		`, idAgence)
		if errBenevole != nil {
			fmt.Printf("Erreur : %v", errBenevole)
			http.Error(response, "Erreur récupération bénévole", http.StatusInternalServerError)
			return
		}
		defer rowsBenevole.Close()

		var benevolesDisponibilites_List []models.TourneeDashboardAffectation

		for rowsBenevole.Next() {
			var benevoleDispo models.TourneeDashboardAffectation

			errScan := rowsBenevole.Scan(&benevoleDispo.ID_Benevole,
				&benevoleDispo.Nom_Benevole,
				&benevoleDispo.Prenom_Benevole,
				&benevoleDispo.ID_Planning,
				&benevoleDispo.Date_Planning,
				&benevoleDispo.Heure_Debut,
				&benevoleDispo.Heure_Fin,
				&benevoleDispo.Statut_Planning,
			)
			if errScan != nil {
				fmt.Printf("Erreur : %v", errScan)
				http.Error(response, "Erreur scan dispo benevole", http.StatusInternalServerError)
				return
			}

			benevolesDisponibilites_List = append(benevolesDisponibilites_List, benevoleDispo)
		}

		data := AdminAgenceBenevolesDisponibilitesTournee{
			Benevole_Disponibilites: benevolesDisponibilites_List,
		}

		tmpl, tmplErr := template.ParseFiles("./templates/admin_agence/benevoles_disponibilites_tournees.html")
		if tmplErr != nil {
			fmt.Printf("Erreur : %v", tmplErr)
			http.Error(response, "Erreur parsefiles html", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "ADMIN_AGENCE")
}
