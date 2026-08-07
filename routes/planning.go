package routes

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
)

const DossierPlannings = "./uploads/plannings"

type LignePlanningExcel struct {
	Date       string
	HeureDebut string
	HeureFin   string
	Intitule   string
	Lieu       string
	Statut     string
}

// CheminPlanningExcel retourne le chemin du fichier excel d'un bénévole donné
func CheminPlanningExcel(idBenevole int) string {
	return filepath.Join(DossierPlannings, fmt.Sprintf("planning_benevole_%d.xlsx", idBenevole))
}

// ouvrirOuCreerClasseur ouvre le classeur existant du bénévole ou en crée un nouveau,
// et s'assure que les 3 feuilles (Services, Collectes, Tournées) existent.
func ouvrirOuCreerClasseur(idBenevole int) (*excelize.File, string, error) {
	chemin := CheminPlanningExcel(idBenevole)

	// même logique que dans inscription.go : os.MkdirAll avant d'écrire un fichier
	if errDir := os.MkdirAll(DossierPlannings, os.ModePerm); errDir != nil {
		return nil, chemin, fmt.Errorf("création dossier plannings : %w", errDir)
	}

	var fichier *excelize.File
	if _, errStat := os.Stat(chemin); errStat == nil {
		f, errOpen := excelize.OpenFile(chemin)
		if errOpen != nil {
			return nil, chemin, fmt.Errorf("ouverture fichier existant : %w", errOpen)
		}
		fichier = f
	} else {
		fichier = excelize.NewFile()
	}

	for _, feuille := range []string{"Services", "Collectes", "Tournées"} {
		if idx, _ := fichier.GetSheetIndex(feuille); idx == -1 {
			if _, errCreate := fichier.NewSheet(feuille); errCreate != nil {
				return nil, chemin, fmt.Errorf("création feuille %s : %w", feuille, errCreate)
			}
		}
	}

	// onglet par défaut d'excelize.NewFile(), à retirer s'il traîne
	if idx, _ := fichier.GetSheetIndex("Sheet1"); idx != -1 {
		fichier.DeleteSheet("Sheet1")
	}

	return fichier, chemin, nil
}

// ecrireFeuillePlanning réécrit entièrement le contenu d'une feuille à partir des lignes fournies
func ecrireFeuillePlanning(fichier *excelize.File, feuille string, lignes []LignePlanningExcel) error {

	fichier.DeleteSheet(feuille)
	if _, errCreate := fichier.NewSheet(feuille); errCreate != nil {
		return errCreate
	}

	entetes := []string{"Date", "Heure début", "Heure fin", "Intitulé", "Lieu", "Statut"}
	for col, entete := range entetes {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		fichier.SetCellValue(feuille, cell, entete)
	}

	styleEntete, errStyle := fichier.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
	})
	if errStyle == nil {
		finCol, _ := excelize.CoordinatesToCellName(len(entetes), 1)
		fichier.SetCellStyle(feuille, "A1", finCol, styleEntete)
	}

	for i, ligne := range lignes {
		numeroLigne := i + 2
		valeurs := []interface{}{ligne.Date, ligne.HeureDebut, ligne.HeureFin, ligne.Intitule, ligne.Lieu, ligne.Statut}
		for col, valeur := range valeurs {
			cell, _ := excelize.CoordinatesToCellName(col+1, numeroLigne)
			fichier.SetCellValue(feuille, cell, valeur)
		}
	}

	for col := range entetes {
		colName, _ := excelize.ColumnNumberToName(col + 1)
		fichier.SetColWidth(feuille, colName, colName, 24)
	}

	return nil
}

// enregistrerDateGenerationPlanning met à jour (ou crée) la ligne planning_export
// du bénévole avec la date/heure actuelle. Appelée à chaque régénération du fichier,
// peu importe l'onglet concerné (Services, Collectes ou Tournées).
func enregistrerDateGenerationPlanning(database *sql.DB, idBenevole int) error {
	_, errExec := database.Exec(`
		INSERT INTO planning_export (id_benevole, date_generation)
		VALUES ($1, NOW())
		ON CONFLICT (id_benevole) DO UPDATE SET date_generation = NOW()
	`, idBenevole)

	if errExec != nil {
		return fmt.Errorf("enregistrement planning_export : %w", errExec)
	}
	return nil
}

// GenererPlanningExcelServices régénère la feuille "Services" du planning excel d'un bénévole
// à partir des demandes de service qui lui sont attribuées. À appeler en fin
// d'AttributionDemandeServicePlanningBenevole.
func GenererPlanningExcelServices(database *sql.DB, idBenevole int) error {

	rows, errQuery := database.Query(`
		SELECT p.date, p.heure_debut, p.heure_fin, s.nom,
		       u.adresse || ', ' || u.ville, ds.statut
		FROM demande_service ds
		JOIN planning p ON p.id_planning = ds.id_planning
		JOIN service s ON s.id_service = ds.id_service
		JOIN adherent a ON a.id_adherent = ds.id_adherent
		JOIN utilisateur u ON u.id_utilisateur = a.id_utilisateur
		WHERE ds.id_benevole = $1
		ORDER BY p.date, p.heure_debut
	`, idBenevole)
	if errQuery != nil {
		return fmt.Errorf("récupération des services attribués : %w", errQuery)
	}
	defer rows.Close()

	var lignes []LignePlanningExcel
	for rows.Next() {
		var date, heureDebut, heureFin time.Time
		var nomService, lieu, statut string
		if errScan := rows.Scan(&date, &heureDebut, &heureFin, &nomService, &lieu, &statut); errScan != nil {
			return fmt.Errorf("scan ligne service : %w", errScan)
		}
		lignes = append(lignes, LignePlanningExcel{
			Date:       date.Format("02/01/2006"),
			HeureDebut: heureDebut.Format("15:04"),
			HeureFin:   heureFin.Format("15:04"),
			Intitule:   nomService,
			Lieu:       lieu,
			Statut:     statut,
		})
	}
	if errRows := rows.Err(); errRows != nil {
		return fmt.Errorf("lecture des lignes services : %w", errRows)
	}

	fichier, chemin, errOuverture := ouvrirOuCreerClasseur(idBenevole)
	if errOuverture != nil {
		return errOuverture
	}

	if errEcriture := ecrireFeuillePlanning(fichier, "Services", lignes); errEcriture != nil {
		return fmt.Errorf("écriture feuille services : %w", errEcriture)
	}

	if errSave := fichier.SaveAs(chemin); errSave != nil {
		return fmt.Errorf("sauvegarde fichier excel : %w", errSave)
	}

	if errDate := enregistrerDateGenerationPlanning(database, idBenevole); errDate != nil {
		return errDate
	}

	return nil
}
