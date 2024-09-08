package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"timur-danilchenko/avito-intership-entry/src/main/database"
	"timur-danilchenko/avito-intership-entry/src/main/models"

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

	var id int
	var createdAt, updatedAt time.Time
	if err = db.QueryRow(query, organization.Name, organization.Description, organization.Type).Scan(&id, &createdAt, &updatedAt); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	log.Infof("Created new organization with ID{%d}", id)

	organization.ID = id
	organization.CreatedAt = createdAt
	organization.UpdatedAt = updatedAt

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

	orgID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

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

	if err := db.QueryRow(query, orgID).Scan(&organization.ID, &organization.Name, &organization.Description, &organization.Type, &organization.CreatedAt, &organization.UpdatedAt); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(organization)
}

func UpdateOrganizationByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	orgID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

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

	if _, err := db.Exec(query, updatedOrganization.Name, updatedOrganization.Description, updatedOrganization.Type, updatedOrganization.UpdatedAt, orgID); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	log.Infof("Updated organization{%d} info", orgID)

	w.WriteHeader(http.StatusNoContent)
}

func DeleteOrganizationByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	orgID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
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
		DELETE FROM organization WHERE id=$1
	`

	if _, err := db.Exec(query, orgID); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	log.Infof("Deleted organization{%d}", orgID)

	w.WriteHeader(http.StatusNoContent)
}

func CreateOrganizationResponsibleHandler(w http.ResponseWriter, r *http.Request) {
	var orgResp models.OrganizationResponsible
	err := json.NewDecoder(r.Body).Decode(&orgResp)
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
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM employee WHERE id=$1)", orgResp.UserID).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "User not exists", http.StatusNotAcceptable)
		return
	}
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM organization WHERE id=$1)", orgResp.OrganizationID).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Organization not exists", http.StatusNotAcceptable)
		return
	}

	query := `INSERT INTO organization_responsible (organization_id, user_id) VALUES ($1, $2) RETURNING id`
	err = db.QueryRow(query, orgResp.OrganizationID, orgResp.UserID).Scan(&orgResp.ID)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Failed to create organization responsible", http.StatusInternalServerError)
		return
	}
	log.Infof("Created new organization responsible with ID{%d}", orgResp.ID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(orgResp)
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

	var orgResps []models.OrganizationResponsible

	for rows.Next() {
		var orgResp models.OrganizationResponsible
		if err := rows.Scan(&orgResp.ID, &orgResp.OrganizationID, &orgResp.UserID); err != nil {
			log.Error(err.Error())
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		orgResps = append(orgResps, orgResp)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orgResps)
}

func GetOrganizationResponsibleByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	orgRespID, err := strconv.Atoi(vars["id"])
	if err != nil {
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

	var orgResp models.OrganizationResponsible
	query := `
		SELECT * FROM organization_responsible WHERE id=$1
	`

	if err := db.QueryRow(query, orgRespID).Scan(&orgResp.ID, &orgResp.OrganizationID, &orgResp.UserID); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orgResp)
}

func UpdateOrganizationResponsibleByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	orgRespID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Invalid organization responsible ID", http.StatusBadRequest)
		return
	}

	var updatedOrgResp models.OrganizationResponsible
	if err := json.NewDecoder(r.Body).Decode(&updatedOrgResp); err != nil {
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

	if _, err := db.Exec(query, updatedOrgResp.OrganizationID, updatedOrgResp.UserID, orgRespID); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	log.Infof("Updated organization responsible{%d} info", orgRespID)
	w.WriteHeader(http.StatusNoContent)
}

func DeleteOrganizationResponsibleByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	orgRespID, err := strconv.Atoi(vars["id"])
	if err != nil {
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

	query := `
		DELETE FROM organization_responsible WHERE id=$1
	`

	if _, err := db.Exec(query, orgRespID); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	log.Infof("Deleted organization responsible{%d}", orgRespID)

	w.WriteHeader(http.StatusNoContent)
}
