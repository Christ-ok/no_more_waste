package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"no_more_waste/models"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func PageInscription(database *sql.DB) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {

		competences, err := getCompetences(database)
		if err != nil {
			fmt.Println("Erreur récupération compétences :", err)
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("templates/inscription.html")
		if err != nil {
			fmt.Println("Erreur parsing template :", err)
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}

		data := struct {
			Competences []models.Competence
		}{
			Competences: competences,
		}

		if err := tmpl.Execute(response, data); err != nil {
			fmt.Println("Erreur exécution template :", err)
		}
	}
}

func getCompetences(database *sql.DB) ([]models.Competence, error) {
	rows, err := database.Query(`SELECT id_competence, nom, description FROM competence ORDER BY nom`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var competences []models.Competence
	for rows.Next() {
		var c models.Competence
		if err := rows.Scan(&c.IDCompetence, &c.Nom, &c.Description); err != nil {
			return nil, err
		}
		competences = append(competences, c)
	}
	return competences, rows.Err()
}

func Signin(database *sql.DB) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {

		if err := request.ParseForm(); err != nil {
			http.Error(response, "Formulaire invalide", http.StatusBadRequest)
			return
		}

		var utilisateur models.Utilisateur
		utilisateur.Nom = request.FormValue("nom")
		utilisateur.Prenom = request.FormValue("prenom")
		utilisateur.Email = request.FormValue("email")
		utilisateur.Telephone = request.FormValue("telephone")
		utilisateur.Adresse = request.FormValue("adresse")
		utilisateur.Ville = request.FormValue("ville")
		utilisateur.CodePostal = request.FormValue("code_postal")
		utilisateur.Pays = request.FormValue("pays")

		motDePasse := request.FormValue("mot_de_passe")
		if motDePasse == "" || motDePasse != request.FormValue("confirmation") {
			http.Error(response, "Les mots de passe ne correspondent pas", http.StatusBadRequest)
			return
		}

		hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(motDePasse), bcrypt.DefaultCost)
		if errHash != nil {
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}
		utilisateur.MotDePasse = string(hashedPassword)

		roleNom := request.FormValue("role")

		var idRole int
		if err := database.QueryRow(
			"SELECT id_role FROM role WHERE nom = $1", roleNom,
		).Scan(&idRole); err != nil {
			http.Error(response, "Rôle invalide", http.StatusBadRequest)
			return
		}

		tx, errTx := database.Begin()
		if errTx != nil {
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		var idUtilisateur int
		errInsert := tx.QueryRow(
			`INSERT INTO utilisateur (nom, prenom, email, mot_de_passe, telephone, adresse, ville, code_postal, pays, id_role)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
			 RETURNING id_utilisateur`,
			utilisateur.Nom, utilisateur.Prenom, utilisateur.Email, utilisateur.MotDePasse,
			utilisateur.Telephone, utilisateur.Adresse, utilisateur.Ville, utilisateur.CodePostal,
			utilisateur.Pays, idRole,
		).Scan(&idUtilisateur)
		if errInsert != nil {
			http.Error(response, "Impossible de créer le compte (email déjà utilisé ?)", http.StatusConflict)
			return
		}

		if err := insertProfilRole(tx, request, roleNom, idUtilisateur); err != nil {
			fmt.Println("Erreur insertion profil rôle :", err)
			http.Error(response, "Impossible d'enregistrer les informations du profil", http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			http.Error(response, "Erreur interne", http.StatusInternalServerError)
			return
		}

		http.Redirect(response, request, "/connexion", http.StatusSeeOther)
	}
}

func insertProfilRole(tx *sql.Tx, request *http.Request, roleNom string, idUtilisateur int) error {
	switch roleNom {

	case "BENEVOLE":
		permis := request.FormValue("permis") == "on"

		var idBenevole int
		err := tx.QueryRow(
			`INSERT INTO benevole (id_utilisateur, permis) VALUES ($1, $2) RETURNING id_benevole`,
			idUtilisateur, permis,
		).Scan(&idBenevole)
		if err != nil {
			return err
		}

		for _, idCompetenceStr := range request.Form["competences"] {
			idCompetence, err := strconv.Atoi(idCompetenceStr)
			if err != nil {
				continue
			}
			if _, err := tx.Exec(
				`INSERT INTO benevole_competence (id_benevole, id_competence) VALUES ($1, $2)`,
				idBenevole, idCompetence,
			); err != nil {
				return err
			}
		}
		return nil

	case "COMMERCANT":
		_, err := tx.Exec(
			`INSERT INTO commercant (id_utilisateur, nom_entreprise, type_commerce, numero_siret)
			 VALUES ($1, $2, $3, $4)`,
			idUtilisateur,
			request.FormValue("nom_entreprise"),
			request.FormValue("type_commerce"),
			request.FormValue("numero_siret"),
		)
		return err

	case "ADHERENT":
		_, err := tx.Exec(`INSERT INTO adherent (id_utilisateur) VALUES ($1)`, idUtilisateur)
		return err

	case "ASSOCIATION":
		nombreBeneficiaires, _ := strconv.Atoi(request.FormValue("nombre_beneficiaires"))
		_, err := tx.Exec(
			`INSERT INTO association_beneficiaire (id_utilisateur, nom_responsable, nom_association, nombre_beneficiaires, type_association)
			 VALUES ($1, $2, $3, $4, $5)`,
			idUtilisateur,
			request.FormValue("nom_responsable"),
			request.FormValue("nom_association"),
			nombreBeneficiaires,
			request.FormValue("type_association"),
		)
		return err

	default:
		return fmt.Errorf("rôle non reconnu : %s", roleNom)
	}
}
