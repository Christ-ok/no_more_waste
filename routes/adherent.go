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

type AdherentData struct {
	Adherent models.Utilisateur
}

type AdhesionPageData struct {
	Statut         string
	DateAdhesion   string
	DateExpiration string
	Succes         string
	Erreur         string
}

func PageAdhesionAdherent(database *sql.DB) http.HandlerFunc {
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

		var statut string
		var dateAdhesion, dateExpiration sql.NullTime

		errRow := database.QueryRow(`
			SELECT statut, date_adhesion, date_expiration
			FROM adherent
			WHERE id_utilisateur = $1
		`, idUtilisateur).Scan(&statut, &dateAdhesion, &dateExpiration)
		if errRow != nil {
			fmt.Printf("Erreur récupération adhérent : %v", errRow)
			http.Error(response, "Erreur récupération adhérent", http.StatusInternalServerError)
			return
		}

		data := AdhesionPageData{
			Statut: statut,
		}

		if dateAdhesion.Valid {
			data.DateAdhesion = dateAdhesion.Time.Format("02/01/2006")
		}
		if dateExpiration.Valid {
			data.DateExpiration = dateExpiration.Time.Format("02/01/2006")
		}

		switch request.URL.Query().Get("succes") {
		case "paiement":
			data.Succes = "Votre cotisation a bien été enregistrée. Merci pour votre soutien !"
		}

		switch request.URL.Query().Get("erreur") {
		case "annule":
			data.Erreur = "Le paiement a été annulé, vous pouvez réessayer à tout moment."
		case "stripe":
			data.Erreur = "Une erreur est survenue lors de la préparation du paiement, veuillez réessayer."
		}

		tmpl, errTmpl := template.ParseFiles("./templates/adherent/adhesion_adherent.html")
		if errTmpl != nil {
			fmt.Printf("Erreur parsing template : %v", errTmpl)
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}

		if errExec := tmpl.Execute(response, data); errExec != nil {
			fmt.Printf("Erreur exécution template : %v", errExec)
		}

	}, "ADHERENT")
}

func AfficherPageModififierProfilAdherent(database *sql.DB) http.HandlerFunc {
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

		rowsErr := database.QueryRow(`SELECT u.nom, u.prenom, u.email, u.telephone, u.adresse, u.ville, u.code_postal, u.pays FROM utilisateur u 
											JOIN adherent a ON u.id_utilisateur = a.id_utilisateur
											WHERE a.id_utilisateur = $1
		`, idUtilisateur).Scan(&user.Nom, &user.Prenom, &user.Email, &user.Telephone, &user.Adresse, &user.Ville, &user.CodePostal, &user.Pays)
		if rowsErr != nil {
			fmt.Printf("Erreur : %v", rowsErr)
			http.Error(response, "Erreur récupération valeurs adhérent", http.StatusInternalServerError)
			return
		}

		data := AdherentData{
			Adherent: user,
		}

		tmpl, tmplErr := template.ParseFiles("./templates/adherent/profile_adherent.html")
		if tmplErr != nil {
			fmt.Printf("Erreur : %v", tmplErr)
			http.Error(response, "Erreur template html", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "ADHERENT")
}

func ModifierProfilAdherent(database *sql.DB) http.HandlerFunc {
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

		http.Redirect(response, request, "/adherent/profil", http.StatusSeeOther)
	}, "ADHERENT")
}

func ModificationMotDePasseAdherent(database *sql.DB) http.HandlerFunc {
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

		http.Redirect(response, request, "/adherent/profil", http.StatusSeeOther)

	}, "ADHERENT")
}

func DashboardAdherentServices(database *sql.DB) http.HandlerFunc {
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

		idAgence, agenceErr := middleware.GetIDAgenceUtilisateur(database, idUtilisateur)
		if agenceErr != nil {
			fmt.Printf("Erreur : %v", agenceErr)
			http.Error(response, "Erreur récupération id agence", http.StatusInternalServerError)
			return
		}

		rowsServices, servicesErr := database.Query(`SELECT id_service, nom, description, statut FROM service WHERE id_agence = $1 AND statut = 'ACTIF'`, idAgence)
		if servicesErr != nil {
			fmt.Printf("Erreur : %v", servicesErr)
			http.Error(response, "Erreur récupération services", http.StatusInternalServerError)
			return
		}
		defer rowsServices.Close()

		var services_List []models.Service

		for rowsServices.Next() {

			var service models.Service

			errQuery := rowsServices.Scan(&service.ID_Service, &service.Nom, &service.Description, &service.Statut)
			if errQuery != nil {
				fmt.Printf("Erreur : %v", errQuery)
				http.Error(response, "Erreur Scan des services", http.StatusInternalServerError)
				return
			}

			services_List = append(services_List, service)
		}

		data := AdminAgenceDashboardServicesAfficher{
			Services: services_List,
		}

		tmpl, tmplErr := template.ParseFiles("./templates/adherent/services.html")
		if tmplErr != nil {
			fmt.Printf("Erreur : %v", tmplErr)
			http.Error(response, "Erreur parsage fichier", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "ADHERENT")
}

func RejoindreServiceAdherent(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		idServiceStr := request.FormValue("id_service")

		idService, err := strconv.Atoi(idServiceStr)
		if err != nil {
			http.Error(response, "ID service invalide", http.StatusBadRequest)
			return
		}

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

		var idAdherent int
		var statut string

		errQuery := database.QueryRow(`SELECT id_adherent, statut FROM adherent WHERE id_utilisateur = $1`, idUtilisateur).Scan(&idAdherent, &statut)
		if errQuery != nil {
			fmt.Printf("Erreur : %v", errQuery)
			http.Error(response, "Erreur récupération ID adhérent", http.StatusInternalServerError)
			return
		}

		if statut != "ACTIF" {
			http.Redirect(response, request, "/adherent/adhesion", http.StatusSeeOther)
			return
		}

		_, errExec := database.Exec(`INSERT INTO demande_service (id_service, id_adherent) VALUES ($1, $2)`, idService, idAdherent)
		if errExec != nil {
			fmt.Printf("Erreur : %v", errExec)
			http.Error(response, "Erreur lors de l'envoie de la demande", http.StatusInternalServerError)
			return
		}

		fmt.Println("Demande de service envoyé avec succès !")
		http.Redirect(response, request, "adherent/services", http.StatusSeeOther)

	}, "ADHERENT")
}

func AfficherServiceRejointAdherent(database *sql.DB) http.HandlerFunc {
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

		var idAdherent int
		errAdherent := database.QueryRow(
			`SELECT id_adherent FROM adherent WHERE id_utilisateur = $1`, idUtilisateur,
		).Scan(&idAdherent)
		if errAdherent != nil {
			fmt.Printf("Erreur : %v", errAdherent)
			http.Error(response, "Erreur récupération adhérent", http.StatusInternalServerError)
			return
		}

		servicesRejoints, serviceErr := database.Query(`SELECT ds.id_demande_service, s.nom, ds.date_demande, ds.statut FROM demande_service ds
														 JOIN service s ON s.id_service = ds.id_service
														WHERE ds.id_adherent = $1
														ORDER BY ds.date_demande DESC
		`, idAdherent)
		if serviceErr != nil {
			fmt.Printf("Erreur : %v", serviceErr)
			http.Error(response, "Erreur récupération service", http.StatusInternalServerError)
			return
		}
		defer servicesRejoints.Close()

		var servicesRejoints_List []models.DemandeService

		for servicesRejoints.Next() {
			var serviceRejoint models.DemandeService

			errScan := servicesRejoints.Scan(&serviceRejoint.ID_Demande_Service,
				&serviceRejoint.Nom_Service,
				&serviceRejoint.Date_Demande,
				&serviceRejoint.Statut)
			if errScan != nil {
				fmt.Printf("Erreur : %v", errScan)
				http.Error(response, "Erreur scan service", http.StatusInternalServerError)
				return
			}

			servicesRejoints_List = append(servicesRejoints_List, serviceRejoint)
		}

		data := AdminAgenceDashboardDemandeService{
			Demandes: servicesRejoints_List,
		}

		tmpl, errTmpl := template.ParseFiles("./templates/adherent/services_rejoints.html")
		if errTmpl != nil {
			fmt.Printf("Erreur : %v", errTmpl)
			http.Error(response, "Erreur parsefile html", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "ADHERENT")
}

func HistoriqueServiceAdherent(database *sql.DB) http.HandlerFunc {
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

		var idAdherent int
		errAdherent := database.QueryRow(
			`SELECT id_adherent FROM adherent WHERE id_utilisateur = $1`, idUtilisateur,
		).Scan(&idAdherent)
		if errAdherent != nil {
			fmt.Printf("Erreur : %v", errAdherent)
			http.Error(response, "Erreur récupération adhérent", http.StatusInternalServerError)
			return
		}

		serviceRendus, errService := database.Query(`SELECT ds.id_demande_service, s.nom, p.date, u.nom, u.prenom, ds.statut FROM demande_service ds
														JOIN service s ON s.id_service = ds.id_service
														JOIN planning p ON p.id_planning = ds.id_planning
														JOIN benevole b ON ds.id_benevole = b.id_benevole
														JOIN utilisateur u ON u.id_utilisateur = b.id_utilisateur
													WHERE ds.id_adherent = $1
														ORDER BY p.date DESC
		`, idAdherent)
		if errService != nil {
			fmt.Printf("Erreur : %v", errService)
			http.Error(response, "Erreur récupération data service", http.StatusInternalServerError)
			return
		}
		defer serviceRendus.Close()

		var serviceRendus_List []models.DemandeService

		for serviceRendus.Next() {
			var serviceRendu models.DemandeService

			errScan := serviceRendus.Scan(&serviceRendu.ID_Demande_Service,
				&serviceRendu.Nom_Service,
				&serviceRendu.Date_Demande,
				&serviceRendu.Nom_Benevole,
				&serviceRendu.Prenom_Benevole,
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

		tmpl, errTmpl := template.ParseFiles("./templates/adherent/historique_services.html")
		if errTmpl != nil {
			fmt.Printf("Erreur : %v", errTmpl)
			http.Error(response, "Erreur parse html file", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(response, data)

	}, "ADHERENT")
}
