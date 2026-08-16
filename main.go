package main

import (
	"fmt"
	"net/http"
	db "no_more_waste/database"
	"no_more_waste/middleware"
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

	if err := os.MkdirAll("./stockage/plannings", 0755); err != nil {
		fmt.Println("Erreur création dossier plannings :", err)
		os.Exit(1)
	}

	http.HandleFunc("GET /inscription", PageInscription(db.DB))
	http.HandleFunc("POST /inscription", Signin(db.DB))
	http.HandleFunc("GET /connexion", PageConnexion)
	http.HandleFunc("POST /connexion", Login(db.DB))
	http.HandleFunc("POST /logout", middleware.Logout)

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

	http.HandleFunc("GET /benevole/disponibilite/creer", routes.PageCreerDisponibilite)
	http.HandleFunc("POST /benevole/disponibilite/creer", routes.BenevoleDisponibilite(db.DB))
	http.HandleFunc("GET /benevole/disponibilites", routes.DashboardBenevoleDisponibiltes(db.DB))
	http.HandleFunc("POST /benevole/disponibilite/modifier", routes.ModifierDisponibilite(db.DB))
	http.HandleFunc("POST /benevole/disponibilite/supprimer", routes.DeleteDisponibilite(db.DB))
	http.HandleFunc("GET /benevole/services/historique", routes.HistoriqueServiceRenduBenevole(db.DB))
	http.HandleFunc("GET /benevole/profil", routes.AfficherPageModifierProfilBenevole(db.DB))
	http.HandleFunc("POST /benevole/profil/modifier", routes.ModifierProfileBenevole(db.DB))
	http.HandleFunc("POST /benevole/profil/mot-de-passe", routes.ModificationMotDePasseBenevole(db.DB))
	http.HandleFunc("GET /benevole/planning-excel", routes.PagePlanningExcelBenevole(db.DB))
	http.HandleFunc("GET /benevole/planning-excel/telecharger", routes.TelechargerPlanningExcel(db.DB))
	http.HandleFunc("GET /benevole/collectes", routes.DashboardBenevoleCollectes(db.DB))
	http.HandleFunc("GET /benevole/tournees", routes.DashboardBenevoleTournees(db.DB))
	http.HandleFunc("POST /benevole/collectes/terminer", routes.TerminerCollecte(db.DB))
	http.HandleFunc("POST /benevole/tournees/terminer", routes.TerminerTournee(db.DB))

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
	http.HandleFunc("GET /admin-agence/collectes", routes.DashboardAdminAgenceCollecte(db.DB))
	http.HandleFunc("GET /admin-agence/demandes-collectes/affectation", routes.PageAffectationBenevoleCollecte(db.DB))
	http.HandleFunc("POST /admin-agence/demandes-collectes/attribuer", routes.AttributionCollectePlanningBenevole(db.DB))
	http.HandleFunc("GET /admin-agence/stocks", routes.DashboardAdminAgenceStocks(db.DB))
	http.HandleFunc("GET /admin-agence/tournees", routes.DashboardAdminAgenceTournee(db.DB))
	http.HandleFunc("GET /admin-agence/tournees/creer", routes.PageCreerTournee(db.DB))
	http.HandleFunc("POST /admin-agence/tournees/creer", routes.CreerTournee(db.DB))
	http.HandleFunc("GET /admin-agence/tournees/affectation", routes.AfficherBenevoleDisponibleTournee(db.DB))
	http.HandleFunc("POST /admin-agence/tournees/attribution", routes.AttributionTourneePlanningBenevole(db.DB))
	http.HandleFunc("GET /admin-agence/recapitulatifs", routes.DashboardAdminAgenceRecapitulatifs(db.DB))
	http.HandleFunc("GET /admin-agence/recapitulatifs/telecharger", routes.TelechargerRecapitulatifTournee(db.DB))

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

	http.HandleFunc("GET /commercant/profil", routes.AfficherPageModifierProfilCommercant(db.DB))
	http.HandleFunc("POST /commercant/profil/modifier", routes.ModifierProfilCommercant(db.DB))
	http.HandleFunc("POST /commercant/profil/mot-de-passe", routes.ModificationMotDePasseCommercant(db.DB))
	http.HandleFunc("GET /commercant/collectes/", routes.DashboardCommercantCollecte(db.DB))
	http.HandleFunc("POST /commercant/collecte/demande", routes.DemandeCollectePourCommercant(db.DB))
	http.HandleFunc("GET /commercant/collecte/demande/page", routes.PageDemandeCollecte)
	http.HandleFunc("GET /commercant/recapitulatifs", routes.DashboardCommercantRecapitulatifs(db.DB))
	http.HandleFunc("GET /commercant/recapitulatifs/telecharger", routes.TelechargerRecapitulatifCollecte(db.DB))

	uploadsFs := http.FileServer(http.Dir("./uploads"))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", uploadsFs))

	fmt.Println("Serveur lancé")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Erreur :", err)
	}
}
