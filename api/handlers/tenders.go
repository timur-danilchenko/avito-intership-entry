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
)

func CreateTenderHandler(w http.ResponseWriter, r *http.Request) {
	var tenderCreate models.TenderCreate
	if err := json.NewDecoder(r.Body).Decode(&tenderCreate); err != nil {
		http.Error(w, fmt.Sprintf("Invalid input: %s", err.Error()), http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	var authorID uuid.UUID
	query := `SELECT id FROM employee WHERE username=$1;`

	if err := db.QueryRow(query, tenderCreate.CreatorUsername).Scan(&authorID); err != nil {
		http.Error(w, fmt.Sprintf("No user with username: {%s}", tenderCreate.CreatorUsername), http.StatusUnauthorized)
		return
	}

	var exists bool
	query = `SELECT EXISTS(SELECT 1 FROM organization WHERE id=$1);`
	if err = db.QueryRow(query, tenderCreate.OrganizationID).Scan(&exists); err != nil || !exists {
		http.Error(w, fmt.Sprintf("No organization with id: {%s}", tenderCreate.OrganizationID), http.StatusUnauthorized)
		return
	}

	query = `SELECT EXISTS(SELECT 1 FROM organization_responsible WHERE organization_id=$1 and user_id=$2);`
	if err := db.QueryRow(query, tenderCreate.OrganizationID, authorID).Scan(&exists); err != nil {
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
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	tender.Name = tenderCreate.Name
	tender.Description = tenderCreate.Description
	tender.ServiceType = tenderCreate.ServiceType
	tender.OrganizationID = tenderCreate.OrganizationID

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
		http.Error(w, "Failed to get tenders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tender models.Tender
		if err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationID, &tender.Version, &tender.CreatedAt); err != nil {
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
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	var userID uuid.UUID
	query := `SELECT id FROM employee WHERE username=$1;`
	if err := db.QueryRow(query, username).Scan(&userID); err != nil {
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
		http.Error(w, "Failed to get tenders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tender models.Tender
		if err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationID, &tender.Version, &tender.CreatedAt); err != nil {
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

	tenderID := vars["tenderId"]
	username := r.URL.Query().Get("username")

	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	var userID uuid.UUID
	query := `SELECT id FROM employee WHERE username=$1;`
	if err := db.QueryRow(query, username).Scan(&userID); err != nil {
		http.Error(w, fmt.Sprintf("No user with username: {%s}", username), http.StatusUnauthorized)
		return
	}

	// TODO: Only organization related users can view tenders

	var tender models.Tender
	query = `
		SELECT status FROM tender WHERE id=$1
	`

	if err := db.QueryRow(query, tenderID).Scan(&tender.Status); err != nil {
		http.Error(w, "Tender not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender.Status)
}

func UpdateTenderStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tenderID := vars["tenderId"]
	if tenderID == "" {
		http.Error(w, "TenderID is required", http.StatusBadRequest)
		return
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		http.Error(w, "Status is required", http.StatusBadRequest)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	var userID uuid.UUID
	query := `SELECT id FROM employee WHERE username=$1;`
	if err := db.QueryRow(query, username).Scan(&userID); err != nil {
		http.Error(w, fmt.Sprintf("No user with username: {%s}", username), http.StatusUnauthorized)
		return
	}

	var organizationID uuid.UUID
	query = `SELECT organization_id FROM tender WHERE id=$1`
	if err := db.QueryRow(query, tenderID).Scan(&organizationID); err != nil {
		http.Error(w, fmt.Sprintf("Tender with TenderID{%s} not found", tenderID), http.StatusNotFound)
		return
	}

	var exists bool
	query = `SELECT EXISTS(SELECT 1 FROM organization_responsible WHERE user_id=$1 AND organization_id=$2);`
	if err := db.QueryRow(query, userID, organizationID).Scan(&exists); err != nil {
		http.Error(w, fmt.Sprintf("User{%s} not responsible for organization{%s}", username, tenderID), http.StatusUnauthorized)
		return
	}

	if !exists {
		http.Error(w, fmt.Sprintf("User{%s} not responsible for organization{%s}", username, tenderID), http.StatusUnauthorized)
		return
	}

	var tender models.Tender
	query = `
		SELECT * FROM tender 
		WHERE id=$1`

	if err := db.QueryRow(query, tenderID).Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationID, &tender.Version, &tender.CreatedAt); err != nil {
		http.Error(w, "Tender not found", http.StatusNotFound)
		return
	}

	if string(tender.Status) == status {
		http.Error(w, "Status is already up to date", http.StatusConflict)
		return
	}

	query = `UPDATE tender SET status=$1 WHERE id=$2`
	if _, err := db.Exec(query, status, tenderID); err != nil {
		http.Error(w, "Failed to update tender status", http.StatusInternalServerError)
		return
	}

	tender.Status = interface{}(status).(models.TenderStatusType)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender)
}

func EditTenderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tenderID := vars["tenderId"]
	if tenderID == "" {
		http.Error(w, "TenderID is required", http.StatusBadRequest)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	var tenderUpdate models.TenderUpdate
	if err := json.NewDecoder(r.Body).Decode(&tenderUpdate); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	var userID uuid.UUID
	query := `SELECT id FROM employee WHERE username=$1`
	if err := db.QueryRow(query, username).Scan(&userID); err != nil {
		http.Error(w, fmt.Sprintf("User with username{%s} not found", username), http.StatusNotFound)
		return
	}

	var organizationID uuid.UUID
	query = `SELECT organization_id FROM tender WHERE id=$1`
	if err := db.QueryRow(query, tenderID).Scan(&organizationID); err != nil {
		http.Error(w, fmt.Sprintf("Tender with TenderID{%s} not found", tenderID), http.StatusNotFound)
		return
	}

	var exists bool
	query = `SELECT EXISTS(SELECT 1 FROM organization_responsible WHERE user_id=$1 AND organization_id=$2);`
	if err := db.QueryRow(query, userID, organizationID).Scan(&exists); err != nil {
		http.Error(w, fmt.Sprintf("User{%s} not responsible for organization{%s}", username, tenderID), http.StatusUnauthorized)
		return
	}

	if !exists {
		http.Error(w, fmt.Sprintf("User{%s} not responsible for organization{%s}", username, tenderID), http.StatusUnauthorized)
		return
	}

	var tender models.Tender
	query = `
		SELECT * FROM tender 
		WHERE id=$1
	`
	if err := db.QueryRow(query, tenderID).Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationID, &tender.Version, &tender.CreatedAt); err != nil {
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

	query = `UPDATE tender SET name=$1, description=$2, service_type=$3 WHERE id=$4`

	if _, err := db.Exec(query, tender.Name, tender.Description, tender.ServiceType, tenderID); err != nil {
		http.Error(w, "Failed to update tender", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender)
}

func RollbackTenderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
}
