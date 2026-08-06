package main

import (
	"fmt"
	"net/http"
	db "no_more_waste/database"
	"no_more_waste/routes"
	"no_more_waste/session"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stripe/stripe-go/v84"
)

func main() {

	godotenv.Load()

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

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
	http.HandleFunc("GET /benevole/disponibilites", routes.DashboardBenevoleDisponibiltes(db.DB))
	http.HandleFunc("POST /benevole/disponibilite/modifier", routes.ModifierDisponibilite(db.DB))
	http.HandleFunc("POST /benevole/disponibilite/supprimer", routes.DeleteDisponibilite(db.DB))

	http.HandleFunc("GET /admin-agence/benevoles", routes.DashboardAdminAgenceBenevoles(db.DB))
	http.HandleFunc("GET /admin-agence/benevoles/disponibilites", routes.DashboardAdminAgenceGererDisponibilite(db.DB))
	http.HandleFunc("GET /admin-agence/benevoles/documents", routes.DashboardAdminAgenceDocumentsBenevoles(db.DB))
	http.HandleFunc("GET /admin-agence/benevoles/documents/voir", routes.VoirDocumentBenevole(db.DB))
	http.HandleFunc("POST /admin-agence/benevoles/valider", routes.ValiderBenevole(db.DB))
	http.HandleFunc("POST /admin-agence/benevoles/rejeter", routes.RejeterBenevole(db.DB))
	http.HandleFunc("GET /admin-agence/planning/creer", routes.FormCreatePlanning(db.DB))
	http.HandleFunc("POST /admin-agence/planning/creer", routes.DashboardAdminAgenceCreerPlanning(db.DB))
	http.HandleFunc("GET /admin-agence/planning", routes.DashboardAdminAgenceAfficherPlanning(db.DB))
	http.HandleFunc("GET /admin-agence/planning/modifier", routes.FormModifierPlanning(db.DB))
	http.HandleFunc("POST /admin-agence/planning/modifier", routes.ModifierPlanning(db.DB))
	http.HandleFunc("DELETE /admin-agence/planning/supprimer", routes.DeletePlanning(db.DB))
	http.HandleFunc("GET /admin-agence/services", routes.DashboardAdminAgenceServices(db.DB))
	http.HandleFunc("GET /admin-agence/services/creer", routes.PageCreerService(db.DB))
	http.HandleFunc("POST /admin-agence/services/creer", routes.CreateService(db.DB))
	http.HandleFunc("GET /admin-agence/services/modifier", routes.FormModifyService(db.DB))
	http.HandleFunc("POST /admin-agence/services/modifier", routes.ModifierService(db.DB))
	http.HandleFunc("DELETE /admin-agence/services/supprimer", routes.DeleteService(db.DB))
	http.HandleFunc("GET /admin-agence/services/historique", routes.HistoriqueServicesRealises(db.DB))
	http.HandleFunc("GET /admin-agence/demande-services", routes.DashboardAdminAgenceDemandesServices(db.DB))
	http.HandleFunc("GET /admin-agence/demandes-services/affectation", routes.PageAffectationBenevoleService(db.DB))
	http.HandleFunc("POST /admin-agence/demandes-services/attribuer", routes.AttributionDemandeServicePlanningBenevole(db.DB))

	http.HandleFunc("GET /adherent/profil", routes.AfficherPageModififierProfilAdherent(db.DB))
	http.HandleFunc("POST /adherent/profil/modifier", routes.ModifierProfilAdherent(db.DB))
	http.HandleFunc("POST /adherent/profil/mot-de-passe", routes.ModificationMotDePasseAdherent(db.DB))
	http.HandleFunc("GET /adherent/adhesion", routes.PageAdhesionAdherent(db.DB))
	http.HandleFunc("POST /adherent/adhesion/payer", routes.CreerSessionPaiementAdhesion(db.DB))
	http.HandleFunc("POST /stripe/webhook", routes.StripeWebhookAdhesion(db.DB, os.Getenv("STRIPE_WEBHOOK_SECRET")))
	http.HandleFunc("GET /adherent/services", routes.DashboardAdherentServices(db.DB))
	http.HandleFunc("POST /adherent/services/rejoindre", routes.RejoindreServiceAdherent(db.DB))
	http.HandleFunc("GET /adherent/services/rejoint", routes.AfficherServiceRejointAdherent(db.DB))
	http.HandleFunc("GET /adherent/services/historique", routes.HistoriqueServiceAdherent(db.DB))

	uploadsFs := http.FileServer(http.Dir("./uploads"))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", uploadsFs))

	fmt.Println("Serveur lancé")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Erreur :", err)
	}
}
