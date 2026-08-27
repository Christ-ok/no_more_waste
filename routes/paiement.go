package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"no_more_waste/middleware"
	"no_more_waste/session"

	"github.com/stripe/stripe-go/v84"
	stripeSession "github.com/stripe/stripe-go/v84/checkout/session"
	"github.com/stripe/stripe-go/v84/webhook"
)

const stripePriceCotisationAdherent = "price_1TzCXjATxxp0L3EN5bLW7Bgf"

func CreerSessionPaiementAdhesion(database *sql.DB) http.HandlerFunc {
	return middleware.AuthRole(func(response http.ResponseWriter, request *http.Request) {

		urlBase := os.Getenv("APP_URL")

		if urlBase == "" {
			fmt.Println("ERREUR : APP_URL EST VIDE")
			http.Error(response, "APP_URL non configurée", http.StatusInternalServerError)
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
		errAdherent := database.QueryRow(
			`SELECT id_adherent FROM adherent WHERE id_utilisateur = $1`, idUtilisateur,
		).Scan(&idAdherent)
		if errAdherent != nil {
			fmt.Printf("Erreur récupération id_adherent : %v", errAdherent)
			http.Error(response, "Erreur récupération adhérent", http.StatusInternalServerError)
			return
		}

		params := &stripe.CheckoutSessionParams{
			PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
			Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),

			LineItems: []*stripe.CheckoutSessionLineItemParams{
				{
					Price:    stripe.String(stripePriceCotisationAdherent),
					Quantity: stripe.Int64(1),
				},
			},

			Metadata: map[string]string{
				"type":        "cotisation_adherent",
				"id_adherent": strconv.Itoa(idAdherent),
			},

			SuccessURL: stripe.String(urlBase + "/adherent/adhesion?succes=paiement"),
			CancelURL:  stripe.String(urlBase + "/adherent/adhesion?erreur=annule"),
		}

		checkoutSess, errStripe := stripeSession.New(params)
		if errStripe != nil {
			fmt.Printf("Erreur création session Stripe : %v", errStripe)
			http.Redirect(response, request, "/adherent/adhesion?erreur=stripe", http.StatusSeeOther)
			return
		}

		http.Redirect(response, request, checkoutSess.URL, http.StatusSeeOther)

	}, "ADHERENT")
}

func StripeWebhookAdhesion(database *sql.DB, webhookSecret string) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {

		fmt.Println(">>> WEBHOOK STRIPE REÇU")
		fmt.Println(">>> METHOD :", request.Method)
		fmt.Println(">>> PATH :", request.URL.Path)

		body, errRead := io.ReadAll(request.Body)
		if errRead != nil {
			http.Error(response, errRead.Error(), http.StatusInternalServerError)
			return
		}

		event, errEvent := webhook.ConstructEvent(
			body,
			request.Header.Get("Stripe-Signature"),
			webhookSecret,
		)
		if errEvent != nil {
			http.Error(response, errEvent.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(">>> SIGNATURE STRIPE VALIDEE :")

		if event.Type != "checkout.session.completed" {
			response.WriteHeader(http.StatusOK)
			return
		}

		fmt.Println(">>> EVENT TYPE :", event.Type)

		var checkoutSess stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &checkoutSess); err != nil {
			http.Error(response, err.Error(), http.StatusBadRequest)
			return
		}

		if checkoutSess.Metadata["type"] != "cotisation_adherent" {
			response.WriteHeader(http.StatusOK)
			return
		}

		fmt.Println(">>> ID ADHERENT :", checkoutSess.Metadata["id_adherent"])
		fmt.Println(">>> TYPE :", checkoutSess.Metadata["type"])
		fmt.Println(">>> SESSION ID :", checkoutSess.ID)
		fmt.Println(">>> TRAITEMENT COTISATION...")

		handleCotisationAdherent(database, checkoutSess)
		response.WriteHeader(http.StatusOK)
	}
}

func handleCotisationAdherent(database *sql.DB, checkoutSess stripe.CheckoutSession) {

	idAdherentStr := checkoutSess.Metadata["id_adherent"]
	if idAdherentStr == "" {
		log.Println("Metadata id_adherent manquante sur la session Stripe")
		return
	}

	idAdherent, errConv := strconv.Atoi(idAdherentStr)
	if errConv != nil {
		log.Printf("Erreur conversion id_adherent : %v", errConv)
		return
	}

	montant := checkoutSess.AmountTotal / 100

	tx, errTx := database.Begin()
	if errTx != nil {
		log.Printf("Erreur ouverture transaction : %v", errTx)
		return
	}
	defer tx.Rollback()

	result, errInsert := tx.Exec(`
		INSERT INTO cotisation (id_adherent, montant, date_paiement, date_debut, date_expiration, statut, stripe_session_id)
		VALUES ($1, $2, NOW(), NOW(), NOW() + INTERVAL '1 year', 'PAYEE', $3)
		ON CONFLICT (stripe_session_id) DO NOTHING
	`, idAdherent, montant, checkoutSess.ID)
	if errInsert != nil {
		log.Printf("Erreur insertion cotisation : %v", errInsert)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("Cotisation déjà enregistrée pour la session %s (webhook rejoué), on ignore", checkoutSess.ID)
		tx.Commit()
		return
	}

	log.Printf(">>> TENTATIVE UPDATE ADHERENT ID %d", idAdherent)

	resultUpdate, errUpdate := tx.Exec(`
		UPDATE adherent
		SET statut          = 'ACTIF',
		    date_adhesion   = COALESCE(date_adhesion, NOW()),
		    date_expiration = NOW() + INTERVAL '1 year'
		WHERE id_adherent = $1
	`, idAdherent)

	log.Printf(">>> UPDATE EXECUTE")

	if errUpdate != nil {
		log.Printf("Erreur mise à jour adhérent : %v", errUpdate)
		return
	}

	rowsUpdated, _ := resultUpdate.RowsAffected()
	log.Printf(">>> NOMBRE D'ADHERENTS MIS A JOUR : %d", rowsUpdated)

	if errCommit := tx.Commit(); errCommit != nil {
		log.Printf("Erreur commit transaction : %v", errCommit)
		return
	}

	log.Printf("Cotisation payée et adhésion activée pour l'adhérent %d (session %s)", idAdherent, checkoutSess.ID)
}
