package models

import (
	"database/sql"
	"time"
)

type Utilisateur struct {
	IDUtilisateur int        `json:"id_utilisateur" db:"id_utilisateur"`
	Nom           string     `json:"nom" db:"nom"`
	Prenom        string     `json:"prenom" db:"prenom"`
	Email         string     `json:"email" db:"email"`
	MotDePasse    string     `json:"-" db:"mot_de_passe"`
	Telephone     string     `json:"telephone" db:"telephone"`
	Adresse       string     `json:"adresse" db:"adresse"`
	Ville         string     `json:"ville" db:"ville"`
	CodePostal    string     `json:"code_postal" db:"code_postal"`
	Pays          string     `json:"pays" db:"pays"`
	DateCreation  time.Time  `json:"date_creation" db:"date_creation"`
	DernierLogin  *time.Time `json:"dernier_login,omitempty" db:"dernier_login"`
	EmailVerifie  bool       `json:"email_verifie" db:"email_verifie"`
	EtatCompte    string     `json:"etat_compte" db:"etat_compte"`
	IDRole        int        `json:"id_role" db:"id_role"`
	IDAgence      *int       `json:"id_agence,omitempty" db:"id_agence"`
}

type Commercant struct {
	IDCommercant   int        `json:"id_commercant" db:"id_commercant"`
	IDUtilisateur  int        `json:"id_utilisateur" db:"id_utilisateur"`
	Nom            string     `json:"nom" db:"nom"`
	Prenom         string     `json:"prenom" db:"prenom"`
	Email          string     `json:"email" db:"email"`
	Adresse        string     `json:"adresse" db:"adresse"`
	NomEntreprise  string     `json:"nom_entreprise" db:"nom_entreprise"`
	TypeCommerce   string     `json:"type_commerce" db:"type_commerce"`
	NumeroSiret    string     `json:"numero_siret" db:"numero_siret"`
	DateAdhesion   *time.Time `json:"date_adhesion,omitempty" db:"date_adhesion"`
	DateExpiration *time.Time `json:"date_expiration,omitempty" db:"date_expiration"`
	Cotisation     float64    `json:"cotisation" db:"cotisation"`
	Statut         string     `json:"statut" db:"statut"`
}

type Benevole struct {
	IDBenevole    int    `json:"id_benevole" db:"id_benevole"`
	IDUtilisateur int    `json:"id_utilisateur" db:"id_utilisateur"`
	Permis        bool   `json:"permis" db:"permis"`
	Disponibilite string `json:"disponibilite" db:"disponibilite"`
	Statut        string `json:"statut" db:"statut"`
}

type Adherent struct {
	IDAdherent     int        `json:"id_adherent" db:"id_adherent"`
	IDUtilisateur  int        `json:"id_utilisateur" db:"id_utilisateur"`
	Nom            string     `json:"nom" db:"nom"`
	Prenom         string     `json:"prenom" db:"prenom"`
	Email          string     `json:"email" db:"email"`
	Adresse        string     `json:"adresse" db:"adresse"`
	DateAdhesion   *time.Time `json:"date_adhesion,omitempty" db:"date_adhesion"`
	DateExpiration *time.Time `json:"date_expiration,omitempty" db:"date_expiration"`
	Cotisation     float64    `json:"cotisation" db:"cotisation"`
	Statut         string     `json:"statut" db:"statut"`
}

type AssociationBeneficiaire struct {
	IDAssociation       int        `json:"id_association" db:"id_association"`
	IDUtilisateur       int        `json:"id_utilisateur" db:"id_utilisateur"`
	NomResponsable      string     `json:"nom_responsable" db:"nom_responsable"`
	NomAssociation      string     `json:"nom_association" db:"nom_association"`
	NombreBeneficiaires int        `json:"nombre_beneficiaires" db:"nombre_beneficiaires"`
	TypeAssociation     string     `json:"type_association" db:"type_association"`
	DateValidation      *time.Time `json:"date_validation,omitempty" db:"date_validation"`
	Statut              string     `json:"statut" db:"statut"`
}

type Competence struct {
	IDCompetence int    `json:"id_competence" db:"id_competence"`
	Nom          string `json:"nom" db:"nom"`
	Description  string `json:"description" db:"description"`
}

type BenevoleCompetence struct {
	IDBenevole   int `json:"id_benevole" db:"id_benevole"`
	IDCompetence int `json:"id_competence" db:"id_competence"`
}

type Role struct {
	IDRole int    `json:"id_role" db:"id_role"`
	Nom    string `json:"nom" db:"nom"`
}

type Agence struct {
	IDAgence   int    `json:"id_agence" db:"id_agence"`
	Nom        string `json:"nom" db:"nom"`
	Adresse    string `json:"adresse" db:"adresse"`
	Ville      string `json:"ville" db:"ville"`
	CodePostal string `json:"code_postal" db:"code_postal"`
	Pays       string `json:"pays" db:"pays"`
	Telephone  string `json:"telephone" db:"telephone"`
	Email      string `json:"email" db:"email"`
}

type BenevoleAffichageDashboard struct {
	IDUtilisateur int    `json:"id_utilisateur" db:"id_utilisateur"`
	Nom           string `json:"nom" db:"nom"`
	Prenom        string `json:"prenom" db:"prenom"`
	Email         string `json:"email" db:"email"`
	Telephone     string `json:"telephone" db:"telephone"`
	Adresse       string `json:"adresse" db:"adresse"`
	Ville         string `json:"ville" db:"ville"`
	CodePostal    string `json:"code_postal" db:"code_postal"`
	Pays          string `json:"pays" db:"pays"`
	Permis        bool   `json:"permis" db:"permis"`
	Disponibilite string `json:"disponibilité" db:"disponibilite"`
	Statut        string `json:"statut" db:"statut"`
}

type CommercantAffichageDashboard struct {
	IDutilisateur int    `json:"id_utilisateur" db:"id_utilisateur"`
	Nom           string `json:"nom" db:"nom"`
	Prenom        string `json:"prenom" db:"prenom"`
	Email         string `json:"email" db:"email"`
	Telephone     string `json:"telephone" db:"telephone"`
	Adresse       string `json:"adresse" db:"adresse"`
	Ville         string `json:"ville" db:"ville"`
	CodePostal    string `json:"code_postal" db:"code_postal"`
	Pays          string `json:"pays" db:"pays"`
	NomEntriprise string `json:"nom_entriprise" db:"nom_entreprise"`
	TypeCommerce  string `json:"type_commerce" db:"type_commerce"`
	NumeroSiret   string `json:"numero_siret" db:"numero_siret"`
	Statut        string `json:"statut" db:"statut"`
}

type AdministrateurAffichageDashboard struct {
	IDUtilisateur int    `json:"id_utilisateur" db:"id_utilisateur"`
	Nom           string `json:"nom" db:"nom"`
	Prenom        string `json:"prenom" db:"prenom"`
	Email         string `json:"email" db:"email"`
	Telephone     string `json:"telephone" db:"telephone"`
	Adresse       string `json:"adresse" db:"adresse"`
	Ville         string `json:"ville" db:"ville"`
	CodePostal    string `json:"code_postal" db:"code_postal"`
	Pays          string `json:"pays" db:"pays"`
	IDAgence      int    `json:"id_agence" db:"id_agence"`
	NomAgence     string `json:"nom_agence" db:"nom_agence"`
	EtatCompte    string `json:"etat_compte" db:"etat_compte"`
}

type BenevoleDisponibilite struct {
	ID_Disponibilite   int       `json:"id_disponibilite" db:"id_disponibilite"`
	ID_Benevole        int       `json:"id_benevole" db:"id_benevole"`
	Nom                string    `json:"nom" db:"nom"`
	Prenom             string    `json:"prenom" db:"prenom"`
	Date_Disponibilite time.Time `json:"date_disponibilite" db:"date_disponibilite"`
	Heure_Debut        time.Time `json:"heure_debut" db:"heure_debut"`
	Heure_Fin          time.Time `json:"heure_fin" db:"heure_fin"`
	Statut             string    `json:"statut" db:"statut"`
	Created_At         time.Time `json:"created_at" db:"created_at"`
}

type Planning struct {
	ID_Planning      int       `json:"id_planning" db:"id_planning"`
	ID_Benevole      int       `json:"id_benevole" db:"id_benevole"`
	ID_Disponibilite int       `json:"id_disponibilite" db:"id_disponibilite"`
	Date             time.Time `json:"date" db:"date"`
	Heure_Debut      time.Time `json:"heure_debut" db:"heure_debut"`
	Heure_Fin        time.Time `json:"heure_fin" db:"heure_fin"`
	Statut           string    `json:"statut" db:"statut"`
}

type PlanningAfficheDashboard struct {
	ID_Planning int       `json:"id_planning" db:"id_planning"`
	Nom         string    `json:"nom" db:"nom"`
	Prenom      string    `json:"prenom" db:"prenom"`
	Date        time.Time `json:"date" db:"date"`
	Heure_Debut time.Time `json:"heure_debut" db:"heure_debut"`
	Heure_Fin   time.Time `json:"heure_fin" db:"heure_fin"`
	Statut      string    `json:"statut" db:"statut"`
}

type Service struct {
	ID_Service  int    `json:"id_service" db:"id_service"`
	Nom         string `json:"nom" db:"nom"`
	Description string `json:"description" db:"description"`
	Competence  string `json:"competence" db:"competence"`
	Statut      string `json:"statut" db:"statut"`
}

type ServiceModifierData struct {
	IDService    int
	Nom          string
	Description  string
	IDCompetence int
	Statut       string
	Competences  []Competence
}

type DemandeService struct {
	ID_Demande_Service int       `json:"id_demande_service" db:"id_demande_service"`
	Nom_Service        string    `json:"nom_service" db:"nom_service"`
	Nom_Adherent       string    `json:"nom_adherent" db:"nom_adherent"`
	Prenom_Adherent    string    `json:"prenom_adherent" db:"prenom_adherent"`
	Nom_Benevole       string    `json:"nom" db:"nom"`
	Prenom_Benevole    string    `json:"prenom" db:"prenom"`
	Date_Demande       time.Time `json:"date_demande" db:"date_demande"`
	Statut             string    `json:"statut" db:"statut"`
}

type BenevoleDisponible struct {
	ID_Benevole int    `json:"id_benevole" db:"id_benevole"`
	Nom         string `json:"nom" db:"nom"`
	Prenom      string `json:"prenom" db:"prenom"`
}

type DemandeServiceDashboard struct {
	ID_Benevole   int       `json:"id_benevole" db:"id_benevole"`
	Nom           string    `json:"nom" db:"nom"`
	Prenom        string    `json:"prenom" db:"prenom"`
	ID_Planning   int       `json:"id_planning" db:"id_planning"`
	Date_Planning time.Time `json:"date_planning" db:"date_planning"`
	Heure_Debut   time.Time `json:"heure_debut" db:"heure_debut"`
	Heure_Fin     time.Time `json:"heure_fin" db:"heure_fin"`
	Statut        string    `json:"statut" db:"statut"`
}

type BenevoleDocument struct {
	ID_Benevole    int    `json:"id_utilisateur" db:"id_utilisateur"`
	Nom            string `json:"nom" db:"nom"`
	Prenom         string `json:"prenom" db:"prenom"`
	Nom_Competence string `json:"nom_competence" db:"nom_competence"`
	Statut         string `json:"statut" db:"statut"`
	Justificatif   string `json:"justificatif" db:"justificatif"`
}

type Planning_Excel struct {
	ID_Planning_Excel int       `json:"id_planning_excel" db:"id_planning_excel"`
	Date_Planning     time.Time `json:"date_planning" db:"date_planning"`
	Heure_Debut       time.Time `json:"heure_debut" db:"heure_debut"`
	Heure_Fin         time.Time `json:"heure_fin" db:"heure_fin"`
	Nom_Competence    string    `json:"nom_competence" db:"nom_competence"`
}

type Collecte struct {
	ID_Collecte       int       `json:"id_collecte" db:"id_collecte"`
	ID_Commercant     int       `json:"id_commercant" db:"id_commercant"`
	Nom_Commercant    string    `json:"nom_commercant" db:"nom_commercant"`
	Prenom_Commercant string    `json:"prenom_commercant" db:"prenom_commercant"`
	ID_Agence         int       `json:"id_agence" db:"id_agence"`
	ID_Benevole       int       `json:"id_benevole" db:"id_benevole"`
	ID_Planning       int       `json:"id_planning" db:"id_planning"`
	Date_Collecte     time.Time `json:"date_collecte" db:"date_collecte"`
	Statut            string    `json:"statut" db:"statut"`
	PDF_Recapitulatif string    `json:"pdf_recapitulatif" db:"pdf_recapitulatif"`
}

type DemandeCollecteDashboard struct {
	ID_Benevole     int       `json:"id_benevole" db:"id_benevole"`
	Nom_Benevole    string    `json:"nom_benevole" db:"nom_benevole"`
	Prenom_Benevole string    `json:"prenom_benevole" db:"prenom_benevole"`
	ID_Planning     int       `json:"id_planning" db:"id_planning"`
	Date_Planning   time.Time `json:"date_planning" db:"date_planning"`
	Heure_Debut     time.Time `json:"heure_debut" db:"heure_debut"`
	Heure_Fin       time.Time `json:"heure_fin" db:"heure_fin"`
	Statut_Planning string    `json:"statut_planning" db:"statut_planning"`
}

type CollecteDashboardBenevole struct {
	ID_Collecte       int       `json:"id_collecte" db:"id_collecte"`
	Nom_Commercant    string    `json:"nom_commercant" db:"nom_commercant"`
	Prenom_Commercant string    `json:"prenom_commercant" db:"prenom_commercant"`
	Date_Collecte     time.Time `json:"date_collecte" db:"date_collecte"`
	Statut_Collecte   string    `json:"statut_collecte" db:"statut_collecte"`
}

type RecapitulatifCollecte struct {
	ID_Collecte   int       `json:"id_collecte" db:"id_collecte"`
	Date_Collecte time.Time `json:"date_collecte" db:"date_collecte"`
}

type Produit struct {
	ID       int
	Quantite float64
}

type StockDashboard struct {
	ID_Stock            int       `json:"id_stock" db:"id_stock"`
	Libelle             string    `json:"libelle" db:"libelle"`
	Code_Barre          string    `json:"code_barre" db:"code_barre"`
	Quantite_Disponible float64   `json:"quantite_disponible" db:"quantite_disponible"`
	Date_Entree         time.Time `json:"date_entree" db:"date_entree"`
}

type TourneeDashboard struct {
	ID_Tournee        int            `json:"id_tournee" db:"id_tournee"`
	Nom_Destinataire  string         `json:"nom_destinataire" db:"nom_destinataire"`
	Type_Destinataire string         `json:"type_destinataire" db:"type_destinataire"`
	Date_Tournee      time.Time      `json:"date_tournee" db:"date_tournee"`
	Statut            string         `json:"statut" db:"statut"`
	Nom_Benevole      sql.NullString `json:"nom_benevole" db:"nom_benevole"`
	Prenom_Benevole   sql.NullString `json:"prenom_benevole" db:"prenom_benevole"`
}

type StockDisponible struct {
	ID_Stock            int     `json:"id_stock" db:"id_stock"`
	Libelle             string  `json:"libelle" db:"libelle"`
	Code_Barre          string  `json:"code_barre" db:"code_barre"`
	Quantite_Disponible float64 `json:"quantite_disponible" db:"quantite_disponible"`
}

type TourneeDashboardBenevole struct {
	ID_Tournee           int       `json:"id_tournee" db:"id_tournee"`
	Nom_Destinataire     string    `json:"nom_destinataire" db:"nom_destinataire"`
	Adresse_Destinataire string    `json:"adresse_destinataire" db:"adresse_destinataire"`
	Date_Tournee         time.Time `json:"date_tournee" db:"date_tournee"`
	Statut_Tournee       string    `json:"statut_tournee" db:"statut_tournee"`
}

type TourneeDashboardAffectation struct {
	ID_Benevole     int       `json:"id_benevole" db:"id_benevole"`
	Nom_Benevole    string    `json:"nom_benevole" db:"nom_benevole"`
	Prenom_Benevole string    `json:"prenom_benevole" db:"prenom_benevole"`
	ID_Planning     int       `json:"id_planning" db:"id_planning"`
	Date_Planning   time.Time `json:"date_planning" db:"date_planning"`
	Heure_Debut     time.Time `json:"heure_debut" db:"heure_debut"`
	Heure_Fin       time.Time `json:"heure_fin" db:"heure_fin"`
	Statut_Planning string    `json:"statut_planning" db:"statut_planning"`
}

type RecapitulatifTournee struct {
	ID_Tournee   int       `json:"id_tournee" db:"id_tournee"`
	Date_Tournee time.Time `json:"date_tournee" db:"date_tournee"`
}
