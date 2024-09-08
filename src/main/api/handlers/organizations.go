package handlers

import (
	"encoding/json"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `
		INSERT INTO organization(name, description, type)
		VALUES($1, $2, $3)
		RETURNING id, created_at, updated_at;
	`

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	var id int
	var createdAt, updatedAt time.Time
	if err = db.QueryRow(sqlStatement, organization.Name, organization.Description, organization.Type).Scan(&id, &createdAt, &updatedAt); err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	sqlStatement := `
		SELECT * FROM organization;
	`

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var organizations []models.Organization

	for rows.Next() {
		var organization models.Organization
		if err := rows.Scan(&organization.ID, &organization.Name, &organization.Description, &organization.Type, &organization.CreatedAt, &organization.UpdatedAt); err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
		log.Fatal(err)
	}
	defer db.Close()

	var organization models.Organization
	sqlStatement := `
		SELECT * FROM organization WHERE id=$1
	`

	if err := db.QueryRow(sqlStatement, orgID).Scan(&organization.ID, &organization.Name, &organization.Description, &organization.Type, &organization.CreatedAt, &organization.UpdatedAt); err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(organization)
}

func UpdateOrganizationByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	orgID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	var updatedOrganization models.Organization
	if err := json.NewDecoder(r.Body).Decode(&updatedOrganization); err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	updatedOrganization.UpdatedAt = time.Now()
	sqlStatement := `
		UPDATE organization SET name=$1, description=$2, type=$3, updated_at=$4 WHERE id=$5
	`

	if _, err := db.Exec(sqlStatement, updatedOrganization.Name, updatedOrganization.Description, updatedOrganization.Type, updatedOrganization.UpdatedAt, orgID); err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Infof("Updated organization{%d} info", orgID)

	w.WriteHeader(http.StatusNoContent)
}
