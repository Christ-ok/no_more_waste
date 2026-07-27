package main

import (
	"fmt"
	"net/http"
	db "no_more_waste/database"
	"no_more_waste/routes"
	"no_more_waste/session"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	godotenv.Load()

	db.Init(".env")
	fmt.Println("Connexion à la base de données réussie")
	db.DB.Exec("SET timezone TO 'Europe/Paris'")

	session.Init()

	fs := http.FileServer(http.Dir("./templates"))
	http.Handle("/", fs)

	http.HandleFunc("GET /inscription", PageInscription(db.DB))
	http.HandleFunc("POST /inscription", Signin(db.DB))
	http.HandleFunc("GET /connexion", PageConnexion)
	http.HandleFunc("POST /connexion", Login(db.DB))

	http.HandleFunc("GET /admin", PageAdminGeneral)
	http.HandleFunc("GET /admin-agence", PageAdminAgence)
	http.HandleFunc("GET /commercant", PageCommercant)
	http.HandleFunc("GET /benevole", PageBenevole)
	http.HandleFunc("GET /association", PageAssociation)
	http.HandleFunc("GET /adherent", PageAdherent)
	http.HandleFunc("GET /admin/benevoles/creer", routes.PageCreerBenevole)
	http.HandleFunc("GET /admin/commercants/creer", routes.PageCreerCommercant)
	http.HandleFunc("GET /benevole/planning", routes.PagePlanningBenevole)
	http.HandleFunc("GET /benevole/dashboard", routes.PageDashboardBenevole)

	http.HandleFunc("GET /admin/benevoles", routes.DashboardAdministrateurBenevoles(db.DB))
	http.HandleFunc("POST /admin/benevoles/creer", routes.CreateBenevole(db.DB))
	http.HandleFunc("GET /admin/benevoles/modifier", routes.FormModifyBenevole(db.DB))
	http.HandleFunc("POST /admin/benevoles/modifier", routes.ModifyBenevole(db.DB))
	http.HandleFunc("DELETE /admin/benevoles/supprimer", routes.DeleteBenevole(db.DB))

	http.HandleFunc("GET /admin/commercants", routes.DashboardAdministrateurCommercant(db.DB))
	http.HandleFunc("POST /admin/commercants/creer", routes.CreateCommercant(db.DB))
	http.HandleFunc("GET /admin/commercants/modifier", routes.FormModifyCommercant(db.DB))
	http.HandleFunc("POST /admin/commercants/modifier", routes.ModifyCommercant(db.DB))
	http.HandleFunc("DELETE /admin/commercants/supprimer", routes.DeleteCommercant(db.DB))

	http.HandleFunc("GET /admin/administrateurs", routes.DashboardAdministrateurAdmins(db.DB))
	http.HandleFunc("GET /admin/administrateurs/creer", routes.PageCreerAdministrateur(db.DB))
	http.HandleFunc("POST /admin/administrateurs/creer", routes.CreateAdministrateur(db.DB))
	http.HandleFunc("GET /admin/administrateurs/modifier", routes.FormModifyAdministrateur(db.DB))
	http.HandleFunc("POST /admin/administrateurs/modifier", routes.ModifyAdministrateur(db.DB))
	http.HandleFunc("DELETE /admin/administrateurs/supprimer", routes.DeleteAdministrateur(db.DB))

	http.HandleFunc("POST /benevole/disponibilite/creer", routes.BenevoleDisponibilite(db.DB))

	http.HandleFunc("GET /admin-agence/benevoles", routes.DashboardAdminAgenceBenevoles(db.DB))
	http.HandleFunc("GET /admin-agence/benevoles/disponibilites", routes.DashboardAdminAgenceGererDisponibilite(db.DB))
	http.HandleFunc("GET /admin-agence/planning/creer", routes.FormCreatePlanning(db.DB))
	http.HandleFunc("POST /admin-agence/planning/creer", routes.DashboardAdminAgenceCreerPlanning(db.DB))
	http.HandleFunc("GET /admin-agence/planning", routes.DashboardAdminAgenceAfficherPlanning(db.DB))
	http.HandleFunc("GET /admin-agence/planning/modifier", routes.FormModifierPlanning(db.DB))
	http.HandleFunc("POST /admin-agence/planning/modifier", routes.ModifierPlanning(db.DB))
	http.HandleFunc("DELETE /admin-agence/planning/supprimer", routes.DeletePlanning(db.DB))

	fmt.Println("Serveur lancé")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Erreur :", err)
	}
}
