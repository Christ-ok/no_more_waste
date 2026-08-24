package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"no_more_waste/i18n"
	"no_more_waste/middleware"
	"no_more_waste/models"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func PageInscription(database *sql.DB) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		language := middleware.GetLanguage(request)

		fmt.Println("LANGUE ACTUELLE :", language)

		competences, err := getCompetences(database)
		if err != nil {
			fmt.Println("Erreur récupération compétences :", err)
			http.Error(response, i18n.Traduction(language, "register.internal_error"), http.StatusInternalServerError)
			return
		}

		agences, agenceErr := getAgences(database)
		if agenceErr != nil {
			fmt.Println("Erreur récupération agences :", err)
			http.Error(response, i18n.Traduction(language, "register.internal_error"), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.New("inscription.html").
			Funcs(template.FuncMap{
				"t": func(key string) string {
					return i18n.Traduction(language, key)
				},
			}).
			ParseFiles("templates/inscription.html")
		if err != nil {
			fmt.Println("Erreur parsing template :", err)
			http.Error(response, i18n.Traduction(language, "register.internal_error"), http.StatusInternalServerError)
			return
		}

		data := struct {
			Competences []models.Competence
			Agences     []models.Agence
		}{
			Competences: competences,
			Agences:     agences,
		}

		if err := tmpl.Execute(response, data); err != nil {
			fmt.Println("Erreur exécution template :", err)
		}
	}
}

func getAgences(database *sql.DB) ([]models.Agence, error) {
	rows, err := database.Query(`SELECT id_agence, nom FROM agence ORDER BY nom`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agences_List []models.Agence

	for rows.Next() {
		var agence models.Agence

		if err := rows.Scan(&agence.IDAgence, &agence.Nom); err != nil {
			return nil, err
		}

		agences_List = append(agences_List, agence)
	}
	return agences_List, rows.Err()
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
		language := middleware.GetLanguage(request)

		fmt.Println("LANGUE ACTUELLE :", language)

		if err := request.ParseMultipartForm(10 << 20); err != nil {
			http.Error(response, i18n.Traduction(language, "register.invalid_form"), http.StatusBadRequest)
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
			http.Error(response, i18n.Traduction(language, "register.password_mismatch"), http.StatusBadRequest)
			return
		}

		hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(motDePasse), bcrypt.DefaultCost)
		if errHash != nil {
			http.Error(response, i18n.Traduction(language, "register.internal_error"), http.StatusInternalServerError)
			return
		}
		utilisateur.MotDePasse = string(hashedPassword)

		roleNom := request.FormValue("role")

		var idRole int
		if err := database.QueryRow(
			"SELECT id_role FROM role WHERE nom = $1", roleNom,
		).Scan(&idRole); err != nil {
			http.Error(response, i18n.Traduction(language, "register.invalid_role"), http.StatusBadRequest)
			return
		}

		idAgence, errAgenceConv := strconv.Atoi(request.FormValue("id_agence"))
		if errAgenceConv != nil {
			http.Error(response, i18n.Traduction(language, "register.invalid_agency"), http.StatusBadRequest)
			return
		}

		tx, errTx := database.Begin()
		if errTx != nil {
			http.Error(response, i18n.Traduction(language, "register.internal_error"), http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		var idUtilisateur int
		errInsert := tx.QueryRow(
			`INSERT INTO utilisateur (nom, prenom, email, mot_de_passe, telephone, adresse, ville, code_postal, pays, id_role, id_agence)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
			 RETURNING id_utilisateur`,
			utilisateur.Nom, utilisateur.Prenom, utilisateur.Email, utilisateur.MotDePasse,
			utilisateur.Telephone, utilisateur.Adresse, utilisateur.Ville, utilisateur.CodePostal,
			utilisateur.Pays, idRole, idAgence,
		).Scan(&idUtilisateur)
		if errInsert != nil {
			http.Error(response, i18n.Traduction(language, "register.account_creation_error"), http.StatusConflict)
			return
		}

		if err := insertProfilRole(tx, request, roleNom, idUtilisateur); err != nil {
			fmt.Println("Erreur insertion profil rôle :", err)
			http.Error(response, i18n.Traduction(language, "register.profile_error"), http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			http.Error(response, i18n.Traduction(language, "register.internal_error"), http.StatusInternalServerError)
			return
		}

		http.Redirect(response, request, "/connexion", http.StatusSeeOther)
	}
}

func insertProfilRole(tx *sql.Tx, request *http.Request, roleNom string, idUtilisateur int) error {
	switch roleNom {

	case "BENEVOLE":
		permis := request.FormValue("permis") == "on"

		idCompetence, errConv := strconv.Atoi(request.FormValue("competence"))
		if errConv != nil {
			return fmt.Errorf("Compétence invalide")
		}

		file, header, errFile := request.FormFile("justificatif")
		if errFile != nil {
			fmt.Printf("Erreur : %v", errFile)
		}
		defer file.Close()

		extension := strings.ToLower(filepath.Ext(header.Filename))
		if extension != ".pdf" {
			return fmt.Errorf("Le justificatid doit etre un fichier PDF")
		}

		if errDir := os.MkdirAll("./uploads", os.ModePerm); errDir != nil {
			return fmt.Errorf("Erreur création dossier uploads : %w", errDir)
		}

		path := filepath.Join("./uploads", fmt.Sprintf("%d.pdf", idUtilisateur))

		out, errCreate := os.Create(path)
		if errCreate != nil {
			return fmt.Errorf("Erreur création fichier : %w", errCreate)
		}
		defer out.Close()

		if _, errCopy := io.Copy(out, file); errCopy != nil {
			return fmt.Errorf("Erreur écriture fichier : %w", errCopy)
		}

		_, err := tx.Exec(`INSERT INTO benevole (id_utilisateur, permis, id_competence, justificatif) VALUES ($1, $2, $3, $4)`, idUtilisateur, permis, idCompetence, path)

		return err

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
		_, err := tx.Exec(`INSERT INTO adherent (id_utilisateur, statut) VALUES ($1, 'EN_ATTENTE')`, idUtilisateur)
		if err != nil {
			fmt.Printf("Erreur : %v", err)
		}
		return err

	default:
		return fmt.Errorf("rôle non reconnu : %s", roleNom)
	}
}
