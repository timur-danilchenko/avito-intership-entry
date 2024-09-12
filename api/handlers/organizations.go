package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"timur-danilchenko/avito-intership-entry/database"
	"timur-danilchenko/avito-intership-entry/models"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func CreateOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	var organization models.Organization
	if err := json.NewDecoder(r.Body).Decode(&organization); err != nil {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("Invalid input: %s", err.Error()), http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO organization(name, description, type)
		VALUES($1, $2, $3)
		RETURNING id, created_at, updated_at;
	`

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if err = db.QueryRow(query, organization.Name, organization.Description, organization.Type).Scan(&organization.ID, &organization.CreatedAt, &organization.UpdatedAt); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	log.Infof("Created new organization with ID{%s}", organization.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(organization)
}

func GetAllOrganizationsHandler(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT * FROM organization;
	`

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Somethig went wrong", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var organizations []models.Organization
	for rows.Next() {
		var organization models.Organization
		if err := rows.Scan(&organization.ID, &organization.Name, &organization.Description, &organization.Type, &organization.CreatedAt, &organization.UpdatedAt); err != nil {
			log.Error(err.Error())
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		organizations = append(organizations, organization)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(organizations)
}

func GetOrganizationByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	organizationID := vars["id"]

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var organization models.Organization
	query := `
		SELECT * FROM organization WHERE id=$1
	`

	if err := db.QueryRow(query, organizationID).Scan(&organization.ID, &organization.Name, &organization.Description, &organization.Type, &organization.CreatedAt, &organization.UpdatedAt); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(organization)
}

func UpdateOrganizationByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	organizationID := vars["id"]

	var updatedOrganization models.Organization
	if err := json.NewDecoder(r.Body).Decode(&updatedOrganization); err != nil {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("Invalid input: %s", err.Error()), http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return

	}
	defer db.Close()

	updatedOrganization.UpdatedAt = time.Now()
	query := `
		UPDATE organization SET name=$1, description=$2, type=$3, updated_at=$4 WHERE id=$5
	`

	if _, err := db.Exec(query, updatedOrganization.Name, updatedOrganization.Description, updatedOrganization.Type, updatedOrganization.UpdatedAt, organizationID); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	log.Infof("Updated organization{%s} info", organizationID)

	w.WriteHeader(http.StatusNoContent)
}

func DeleteOrganizationByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	organizationID := vars["id"]

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `
		DELETE FROM organization WHERE id=$1
	`

	if _, err := db.Exec(query, organizationID); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	log.Infof("Deleted organization{%s}", organizationID)

	w.WriteHeader(http.StatusNoContent)
}

func CreateOrganizationResponsibleHandler(w http.ResponseWriter, r *http.Request) {
	var organizationResponsible models.OrganizationResponsible
	err := json.NewDecoder(r.Body).Decode(&organizationResponsible)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("Invalid input: %s", err.Error()), http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM employee WHERE id=$1)", organizationResponsible.UserID).Scan(&exists)
	if err != nil || !exists {
		log.Error(err.Error())
		http.Error(w, "User not exists", http.StatusNotAcceptable)
		return
	}
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM organization WHERE id=$1)", organizationResponsible.OrganizationID).Scan(&exists)
	if err != nil || !exists {
		log.Error(err.Error())
		http.Error(w, "Organization not exists", http.StatusNotAcceptable)
		return
	}

	query := `INSERT INTO organization_responsible (organization_id, user_id) VALUES ($1, $2) RETURNING id`
	err = db.QueryRow(query, organizationResponsible.OrganizationID, organizationResponsible.UserID).Scan(&organizationResponsible.ID)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Failed to create organization responsible", http.StatusInternalServerError)
		return
	}
	log.Infof("Created new organization responsible with ID{%s}", organizationResponsible.ID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(organizationResponsible)
}

func GetAllOrganizationsResponsiblesHandler(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT * FROM organization_responsible;
	`

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Somethig went wrong", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var organizationResponsibles []models.OrganizationResponsible
	for rows.Next() {
		var organizationResponsible models.OrganizationResponsible
		if err := rows.Scan(&organizationResponsible.ID, &organizationResponsible.OrganizationID, &organizationResponsible.UserID); err != nil {
			log.Error(err.Error())
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		organizationResponsibles = append(organizationResponsibles, organizationResponsible)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(organizationResponsibles)
}

func GetOrganizationResponsibleByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	organizationResponsibleID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Invalid organization responsible ID", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var organizationResponsible models.OrganizationResponsible
	query := `
		SELECT * FROM organization_responsible WHERE id=$1
	`

	if err := db.QueryRow(query, organizationResponsibleID).Scan(&organizationResponsible.ID, &organizationResponsible.OrganizationID, &organizationResponsible.UserID); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(organizationResponsible)
}

func UpdateOrganizationResponsibleByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	organizationResponsibleID := vars["id"]

	var updatedOrganizationResponsible models.OrganizationResponsible
	if err := json.NewDecoder(r.Body).Decode(&updatedOrganizationResponsible); err != nil {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("Invalid input: %s", err.Error()), http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return

	}
	defer db.Close()

	query := `
		UPDATE organization_responsible SET organization_id=$1, user_id=$2 WHERE id=$3
	`

	if _, err := db.Exec(query, updatedOrganizationResponsible.OrganizationID, updatedOrganizationResponsible.UserID, organizationResponsibleID); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	log.Infof("Updated organization responsible{%s} info", organizationResponsibleID)
	w.WriteHeader(http.StatusNoContent)
}

func DeleteOrganizationResponsibleByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	organizationResponsibleID := vars["id"]

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `
		DELETE FROM organization_responsible WHERE id=$1
	`

	if _, err := db.Exec(query, organizationResponsibleID); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	log.Infof("Deleted organization responsible{%s}", organizationResponsibleID)

	w.WriteHeader(http.StatusNoContent)
}
