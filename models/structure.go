package models

import (
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
