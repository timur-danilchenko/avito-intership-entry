package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"timur-danilchenko/avito-intership-entry/database"
	"timur-danilchenko/avito-intership-entry/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func CreateTenderHandler(w http.ResponseWriter, r *http.Request) {
	var tenderCreate models.TenderCreate
	if err := json.NewDecoder(r.Body).Decode(&tenderCreate); err != nil {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("Invalid input: %s", err.Error()), http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	var authorID uuid.UUID
	query := `SELECT id FROM employee WHERE username=$1;`

	if err := db.QueryRow(query, tenderCreate.CreatorUsername).Scan(&authorID); err != nil {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("No user with username: {%s}", tenderCreate.CreatorUsername), http.StatusUnauthorized)
		return
	}

	var exists bool
	query = `SELECT EXISTS(SELECT 1 FROM organization WHERE id=$1);`
	if err = db.QueryRow(query, tenderCreate.OrganizationID).Scan(&exists); err != nil || !exists {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("No organization with id: {%s}", tenderCreate.OrganizationID), http.StatusUnauthorized)
		return
	}

	query = `SELECT EXISTS(SELECT 1 FROM organization_responsible WHERE organization_id=$1 and user_id=$2);`
	if err := db.QueryRow(query, tenderCreate.OrganizationID, authorID).Scan(&exists); err != nil {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("Username with id{%s} are not responsible", tenderCreate.CreatorUsername), http.StatusForbidden)
		return
	}

	query = `
		INSERT INTO tender(name, description, service_type, organization_id)
		VALUES($1, $2, $3, $4)
		RETURNING id, status, version, created_at;
	`

	var tender models.Tender
	if err = db.QueryRow(query, tenderCreate.Name, tenderCreate.Description, tenderCreate.ServiceType, tenderCreate.OrganizationID).Scan(&tender.ID, &tender.Status, &tender.Version, &tender.CreatedAt); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	tender.Name = tenderCreate.Name
	tender.Description = tenderCreate.Description
	tender.ServiceType = tenderCreate.ServiceType
	tender.OrganizationID = tenderCreate.OrganizationID

	log.Infof("Created new tender with ID{%s}", tender.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tender)
}

func GetAllTendersHandler(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 5
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	serviceTypes := r.URL.Query()["service_type"]

	// TODO: Add enum validate

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	var tenders []models.Tender
	query := `
		SELECT * FROM tender
	`

	if len(serviceTypes) > 0 {
		query += " WHERE service_type IN ("
		for i, serviceType := range serviceTypes {
			if i > 0 {
				query += ", "
			}
			query += "'" + serviceType + "'"
		}
		query += ")"
	}

	query += " ORDER BY name LIMIT $1 OFFSET $2"

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Failed to get tenders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tender models.Tender
		if err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationID, &tender.Version, &tender.CreatedAt); err != nil {
			log.Error(err.Error())
			http.Error(w, "Failed to get tenders", http.StatusInternalServerError)
			return
		}
		tenders = append(tenders, tender)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenders)
}

func GetUserTendersHandler(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 5
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		log.Error("Username is required")
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	var userID uuid.UUID
	query := `SELECT id FROM employee WHERE username=$1;`
	if err := db.QueryRow(query, username).Scan(&userID); err != nil {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("No user with username: {%s}", username), http.StatusUnauthorized)
		return
	}

	var tenders []models.Tender
	query = `
		SELECT t.* FROM tender t
		JOIN organization o ON t.organization_id = o.id
		JOIN organization_responsible op ON op.organization_id = o.id
		JOIN employee e ON op.user_id = e.id
		WHERE e.username = $1
		ORDER BY t.name LIMIT $2 OFFSET $3
	`

	rows, err := db.Query(query, username, limit, offset)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Failed to get tenders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tender models.Tender
		if err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationID, &tender.Version, &tender.CreatedAt); err != nil {
			log.Error(err.Error())
			http.Error(w, "Failed to get tenders", http.StatusInternalServerError)
			return
		}
		tenders = append(tenders, tender)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenders)
}

func GetTenderStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tenderID, err := strconv.Atoi(vars["tenderId"])
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Invalid tender ID", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	var tender models.Tender
	query := `
		SELECT status FROM tender WHERE id=$1
	`

	if err := db.QueryRow(query, tenderID).Scan(&tender.Status); err != nil {
		log.Error(err.Error())
		http.Error(w, "Tender not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender.Status)
}

func UpdateTenderStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tenderID, err := strconv.Atoi(vars["tenderId"])
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Invalid tender ID", http.StatusBadRequest)
		return
	}

	status := vars["status"]
	if status == "" {
		log.Error("Status is required")
		http.Error(w, "Status is required", http.StatusBadRequest)
		return
	}

	username := vars["username"]
	if username == "" {
		log.Error("Username is required")
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	var tender models.Tender
	query := `
		SELECT * FROM tender WHERE id=$1
	`

	if err := db.QueryRow(query, tenderID).Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationID, &tender.Version, &tender.CreatedAt); err != nil {
		log.Error(err.Error())
		http.Error(w, "Tender not found", http.StatusNotFound)
		return
	}

	if tender.Status == status {
		log.Info("Status is already up to date")
		http.Error(w, "Status is already up to date", http.StatusConflict)
		return
	}

	query = `
		UPDATE tender SET status=$1 WHERE id=$2
	`

	if _, err := db.Exec(query, status, tenderID); err != nil {
		log.Error(err.Error())
		http.Error(w, "Failed to update tender status", http.StatusInternalServerError)
		return
	}

	tender.Status = status

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender)
}

func EditTenderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tenderID, err := strconv.Atoi(vars["tenderId"])
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Invalid tender ID", http.StatusBadRequest)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		log.Error("Username is required")
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	var tenderUpdate models.TenderUpdate
	if err := json.NewDecoder(r.Body).Decode(&tenderUpdate); err != nil {
		log.Error(err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	var tender models.Tender
	query := `
		SELECT * FROM tender WHERE id=$1
	`

	if err := db.QueryRow(query, tenderID).Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationID, &tender.Version, &tender.CreatedAt); err != nil {
		log.Error(err.Error())
		http.Error(w, "Tender not found", http.StatusNotFound)
		return
	}

	if tenderUpdate.Name != "" {
		tender.Name = tenderUpdate.Name
	}

	if tenderUpdate.Description != "" {
		tender.Description = tenderUpdate.Description
	}

	if tenderUpdate.ServiceType != "" {
		tender.ServiceType = tenderUpdate.ServiceType
	}

	query = `
		UPDATE tender SET name=$1, description=$2, service_type=$3 WHERE id=$4
	`

	if _, err := db.Exec(query, tender.Name, tender.Description, tender.ServiceType, tenderID); err != nil {
		log.Error(err.Error())
		http.Error(w, "Failed to update tender", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender)
}

func RollbackTenderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tenderID, err := strconv.Atoi(vars["tenderId"])
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Invalid tender ID", http.StatusBadRequest)
		return
	}

	version, err := strconv.Atoi(vars["version"])
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		log.Error("Username is required")
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	var tender models.Tender
	query := `
		SELECT * FROM tender WHERE id=$1
	`

	if err := db.QueryRow(query, tenderID).Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationID, &tender.Version, &tender.CreatedAt); err != nil {
		log.Error(err.Error())
		http.Error(w, "Tender not found", http.StatusNotFound)
		return
	}

	if version > tender.Version {
		log.Error("Version is too high")
		http.Error(w, "Version is too high", http.StatusBadRequest)
		return
	}

	query = `
		SELECT name, description, service_type FROM tender_history WHERE tender_id=$1 AND version=$2
	`

	var name, description, serviceType string
	if err := db.QueryRow(query, tenderID, version).Scan(&name, &description, &serviceType); err != nil {
		log.Error(err.Error())
		http.Error(w, "Version not found", http.StatusNotFound)
		return
	}

	tender.Name = name
	tender.Description = description
	tender.ServiceType = serviceType
	tender.Version++

	query = `
		UPDATE tender SET name=$1, description=$2, service_type=$3, version=$4 WHERE id=$5
	`

	if _, err := db.Exec(query, tender.Name, tender.Description, tender.ServiceType, tender.Version, tenderID); err != nil {
		log.Error(err.Error())
		http.Error(w, "Failed to update tender", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender)
}
