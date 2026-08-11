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

type BenevoleDashboardDisponibilite struct {
	Disponibilites []models.BenevoleDisponibilite
}

type BenevoleData struct {
	Benevole   models.Utilisateur
	Competence string
}

type BenevolePlanningExcel struct {
	Planning []models.Planning_Excel
}

func PageDashboardBenevole(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "templates/benevole/accueil_Benevole.html")
}

func PagePlanningBenevole(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "templates/benevole/planning.html")
}

func PageCreerDisponibilite(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "templates/benevole/creer_disponibilite.html")
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
		http.Redirect(response, request, "/benevole/disponibilites", http.StatusSeeOther)

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
																WHERE id_benevole = $1 AND (date_disponibilite + heure_fin) >= NOW () 
																ORDER BY date_disponibilite, heure_debut`,
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

func HistoriqueServiceRenduBenevole(database *sql.DB) http.HandlerFunc {
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

		rowsServices, errServices := database.Query(`SELECT ds.id_demande_service, s.nom, p.date, u.nom, u.prenom, ds.statut FROM demande_service ds
													   JOIN service s ON s.id_service = ds.id_service
													   JOIN planning p ON p.id_planning = ds.id_planning
													   JOIN adherent a ON ds.id_adherent = a.id_adherent
													   JOIN utilisateur u ON u.id_utilisateur = a.id_utilisateur
													 WHERE ds.id_benevole = $1
		`, idBenevole)
		if errServices != nil {
			fmt.Printf("Erreur : %v", errServices)
			http.Error(response, "Erreur récupération de services", http.StatusInternalServerError)
			return
		}
		defer rowsServices.Close()

		var serviceRendus_List []models.DemandeService

		for rowsServices.Next() {
			var serviceRendu models.DemandeService

			errScan := rowsServices.Scan(&serviceRendu.ID_Demande_Service,
				&serviceRendu.Nom_Service,
				&serviceRendu.Date_Demande,
				&serviceRendu.Nom_Adherent,
				&serviceRendu.Prenom_Adherent,
				&serviceRendu.Statut)
			if errScan != nil {
				fmt.Printf("Erreur : %v", errScan)
				http.Error(response, "Erreur Scan service rendu", http.StatusInternalServerError)
				return
			}

			serviceRendus_List = append(serviceRendus_List, serviceRendu)
		}

		data := AdminAgenceDashboardDemandeService{
			Demandes: serviceRendus_List,
		}

		tmpl, errTmpl := template.ParseFiles("./templates/benevole/historique_services.html")
		if errTmpl != nil {
			fmt.Printf("Erreur : %v", errTmpl)
			http.Error(response, "Erreur parsefile html", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "BENEVOLE")
}

func AfficherPageModifierProfilBenevole(database *sql.DB) http.HandlerFunc {
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
		var competence_nom string

		rowsErr := database.QueryRow(`SELECT u.nom, u.prenom, u.email, u.telephone, c.nom, u.ville, u.code_postal, u.pays FROM utilisateur u
									  JOIN benevole b ON b.id_utilisateur = u.id_utilisateur
									  JOIN competence c ON c.id_competence = b.id_competence
									  WHERE b.id_utilisateur = $1
		`, idUtilisateur).Scan(&user.Nom, &user.Prenom, &user.Email, &user.Telephone, &competence_nom, &user.Ville, &user.CodePostal, &user.Pays)
		if rowsErr != nil {
			fmt.Printf("Erreur : %v", rowsErr)
			http.Error(response, "Erreur récupération valeurs bénévole", http.StatusInternalServerError)
			return
		}

		data := BenevoleData{
			Benevole:   user,
			Competence: competence_nom,
		}

		tmpl, errTmpl := template.ParseFiles("./templates/benevole/profil_benevole.html")
		if errTmpl != nil {
			fmt.Printf("Erreur : %v", errTmpl)
			http.Error(response, "Erreur prsefiles html", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "BENEVOLE")
}

func ModifierProfileBenevole(database *sql.DB) http.HandlerFunc {
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
		ville := request.FormValue("ville")
		codePostal := request.FormValue("code_postal")
		pays := request.FormValue("pays")

		_, errExec := database.Exec(`UPDATE utilisateur SET 
						nom         = COALESCE(NULLIF($1, ''), nom),
						prenom      = COALESCE(NULLIF($2, ''), prenom),
						email       = COALESCE(NULLIF($3, ''), email),
						telephone   = COALESCE(NULLIF($4, ''), telephone),
						ville       = COALESCE(NULLIF($5, ''), ville),
						code_postal = COALESCE(NULLIF($6, ''), code_postal),
						pays        = COALESCE(NULLIF($7, ''), pays)
					WHERE id_utilisateur = $8									
		`, nom, prenom, email, telephone, ville, codePostal, pays, idUtilisateur)
		if errExec != nil {
			fmt.Printf("Erreur : %v", errExec)
			http.Error(response, "Erreur Update data utilisateur", http.StatusInternalServerError)
			return
		}

		http.Redirect(response, request, "/benevole/profil", http.StatusSeeOther)

	}, "BENEVOLE")
}

func ModificationMotDePasseBenevole(database *sql.DB) http.HandlerFunc {
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

		http.Redirect(response, request, "/benevole/profil", http.StatusSeeOther)

	}, "BENEVOLE")
}

func PagePlanningExcelBenevole(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, _ := session.Store.Get(request, "nmw-session")
		idUtilisateur := sess.Values["id_utilisateur"].(int)

		var idBenevole int
		err := database.QueryRow(`SELECT id_benevole FROM benevole WHERE id_utilisateur = $1`, idUtilisateur).Scan(&idBenevole)
		if err != nil {
			http.Error(response, "Bénévole introuvable", http.StatusForbidden)
			return
		}

		rows, err := database.Query(`
			SELECT pe.id_planning_excel, p.date, p.heure_debut, p.heure_fin, s.nom
			FROM planning_excel pe
			JOIN planning p ON p.id_planning = pe.id_planning
			JOIN demande_service ds ON ds.id_planning = p.id_planning
			JOIN service s ON s.id_service = ds.id_service
			WHERE pe.id_benevole = $1
			ORDER BY p.date
		`, idBenevole)
		if err != nil {
			fmt.Printf("Erreur : %v", err)
			http.Error(response, "Erreur récupération données excel", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var plannings_List []models.Planning_Excel

		for rows.Next() {
			var planning models.Planning_Excel

			errScan := rows.Scan(&planning.ID_Planning_Excel,
				&planning.Date_Planning,
				&planning.Heure_Debut,
				&planning.Heure_Fin,
				&planning.Nom_Competence,
			)
			if errScan != nil {
				fmt.Printf("Erreur : %v", errScan)
				http.Error(response, "Erreur scan planning", http.StatusInternalServerError)
				return
			}

			plannings_List = append(plannings_List, planning)
		}

		data := BenevolePlanningExcel{
			Planning: plannings_List,
		}

		tmpl, errTmpl := template.ParseFiles("./templates/benevole/planning_excel.html")
		if errTmpl != nil {
			fmt.Printf("Erreur : %v", errTmpl)
			http.Error(response, "Erreur parsefile html", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "BENEVOLE")
}

func TelechargerPlanningExcel(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		sess, sessErr := session.Store.Get(request, "nmw-session")
		if sessErr != nil {
			http.Error(response, "Erreur récupération session", http.StatusInternalServerError)
			return
		}
		idUtilisateur := sess.Values["id_utilisateur"].(int)

		if errForm := request.ParseForm(); errForm != nil {
			http.Error(response, "Erreur formulaire", http.StatusBadRequest)
			return
		}

		idPlanningExcel, err := strconv.Atoi(request.URL.Query().Get("id"))
		if err != nil {
			http.Error(response, "ID invalide", http.StatusBadRequest)
			return
		}

		var cheminFichier string
		err = database.QueryRow(`
			SELECT pe.chemin_fichier
			FROM planning_excel pe
			JOIN benevole b ON b.id_benevole = pe.id_benevole
			WHERE pe.id_planning_excel = $1 AND b.id_utilisateur = $2
		`, idPlanningExcel, idUtilisateur).Scan(&cheminFichier)

		if err != nil {
			http.Error(response, "Fichier introuvable", http.StatusNotFound)
			return
		}

		response.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		response.Header().Set("Content-Disposition", `attachment; filename="planning.xlsx"`)
		http.ServeFile(response, request, cheminFichier)

	}, "BENEVOLE")
}
