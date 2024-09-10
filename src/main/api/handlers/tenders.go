package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"timur-danilchenko/avito-intership-entry/src/main/database"
	"timur-danilchenko/avito-intership-entry/src/main/models"

	log "github.com/sirupsen/logrus"
)

func CreateTenderHandler(w http.ResponseWriter, r *http.Request) {
	var tender models.Tender
	if err := json.NewDecoder(r.Body).Decode(&tender); err != nil {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("Invalid input: %s", err.Error()), http.StatusBadRequest)
		return
	}

	log.Println(tender.ServiceType)

	query := `
		INSERT INTO tenders(name, description, service_type, status, organization_id)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id, created_at;
	`

	// name, description, serviceType, status, organizationId, creatorUsername

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	if err = db.QueryRow(query, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationID).Scan(&tender.ID, &tender.CreatedAt); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	log.Infof("Created new tender with ID{%d}", tender.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tender)
}
