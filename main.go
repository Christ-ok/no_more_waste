package main

import (
	"fmt"
	"log"
	db "no_more_waste/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	godotenv.Load()

	db.Init(".env")
	fmt.Println("Connexion à la base de données réussie")

	db.DB.Exec("SET timezone TO 'Europe/Paris'")

	rows, err := db.DB.Query("SELECT id_role, nom FROM role")
	if err != nil {
		log.Fatalf("Erreur lors de la requête: %v", err)
	}
	defer rows.Close()

	fmt.Println("Voici la table role :")
	for rows.Next() {
		var idRole int
		var nom string

		if err := rows.Scan(&idRole, &nom); err != nil {
			log.Fatalf("Erreur lors du scan: %v", err)
		}

		fmt.Printf("id_role: %d | nom: %s\n", idRole, nom)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Erreur après itération: %v", err)
	}
}
