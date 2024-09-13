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

func CreateBidHandler(w http.ResponseWriter, r *http.Request) {
	var bidCreate models.BidCreate
	if err := json.NewDecoder(r.Body).Decode(&bidCreate); err != nil {
		http.Error(w, fmt.Sprintf("Invalid input: %s", err.Error()), http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	var status string
	query := `SELECT status FROM tender WHERE id=$1;`
	if err := db.QueryRow(query, bidCreate.TenderID).Scan(&status); err != nil {
		http.Error(w, fmt.Sprintf("No tender with id: {%s}", bidCreate.TenderID), http.StatusNotFound)
		return
	}

	if status != "Published" {
		http.Error(w, fmt.Sprintf("Tender{%s} status isn't published", bidCreate.TenderID), http.StatusNotFound)
		return
	}

	var bid models.Bid
	query = `
		INSERT INTO bid (name, description, tender_id, author_type, author_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, version, created_at;
	`

	if err := db.QueryRow(query, bidCreate.Name, bidCreate.Description, bidCreate.TenderID, bidCreate.AuthorType, bidCreate.AuthorID).Scan(&bid.ID, &bid.Version, &bid.CreatedAt); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	bid.Name = bidCreate.Name
	bid.Description = bidCreate.Description
	bid.TenderID = bidCreate.TenderID
	bid.AuthorType = bidCreate.AuthorType
	bid.AuthorID = bidCreate.AuthorID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bid)
}

func GetUserBidsHandler(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
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
	query := `SELECT id FROM employee WHERE username=$1`
	if err := db.QueryRow(query, username).Scan(&userID); err != nil {
		http.Error(w, fmt.Sprintf("No user found with username{%s}", username), http.StatusServiceUnavailable)
		return
	}

	var organizationID uuid.UUID
	query = `SELECT organization_id FROM organization_responsible WHERE user_id=$1`
	if err = db.QueryRow(query, userID).Scan(&organizationID); err != nil {
		http.Error(w, fmt.Sprintf("No user found with username{%s}", username), http.StatusServiceUnavailable)
		return
	}

	var bids []models.Bid
	query = `
		SELECT * FROM bid 
		WHERE author_id=$1
		ORDER BY name LIMIT $2 OFFSET $3
	`

	rows, err := db.Query(query, userID, limit, offset)
	if err != nil {
		http.Error(w, "Failed to get bids", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var bid models.Bid
		if err := rows.Scan(&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID, &bid.AuthorType, &bid.AuthorID, &bid.Version, &bid.CreatedAt); err != nil {
			http.Error(w, "Failed to get bids", http.StatusInternalServerError)
			return
		}
		bids = append(bids, bid)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bids)
}

func GetBidsForTenderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderID := vars["tenderId"]

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 5
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	query := `
		SELECT * FROM bid
		WHERE tender_id=$1
		ORDER BY name LIMIT $2 OFFSET $3
	`
	var bids []models.Bid

	rows, err := db.Query(query, tenderID, limit, offset)
	if err != nil {
		http.Error(w, "Can't get tenders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var bid models.Bid
		if err := rows.Scan(&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID, &bid.AuthorType, &bid.AuthorID, &bid.Version, &bid.CreatedAt); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		bids = append(bids, bid)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bids)
}

func GetBidStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID := vars["bidId"]

	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username не указан", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	var bidAuthorType models.BidAuthorType
	query := `SELECT author_type FROM bid WHERE id=$1`
	if err := db.QueryRow(query, bidID).Scan(&bidAuthorType); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var authorID uuid.UUID
	if bidAuthorType == models.BidAuthorOrganization {
		query = `
			SELECT o.id FROM organization o
			JOIN organization_responsible orgre ON orgre.organization_id=o.id
			JOIN employee e ON orgre.user_id=e.id
			WHERE e.username=$1;
		
		`
	} else if bidAuthorType == models.BidAuthorUser {
		query = `SELECT id FROM employee WHERE username=$1`
	}

	if err = db.QueryRow(query, username).Scan(&authorID); err != nil {
		http.Error(w, "Can't find author id", http.StatusNotFound)
		return
	}

	query = `
		SELECT status FROM bid
		WHERE id=$1 AND author_id=$2
	`
	var status string
	if err := db.QueryRow(query, bidID, authorID).Scan(&status); err != nil {
		http.Error(w, "AuthorID incorrect or not responsible organization", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

func UpdateBidStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID := vars["bidId"]

	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username не указан", http.StatusBadRequest)
		return
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		http.Error(w, "status не указан", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	var bidAuthorType models.BidAuthorType
	query := `SELECT author_type FROM bid WHERE id=$1`
	if err := db.QueryRow(query, bidID).Scan(&bidAuthorType); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var authorID uuid.UUID
	if bidAuthorType == models.BidAuthorOrganization {
		query = `
			SELECT o.id FROM organization o
			JOIN organization_responsible orgre ON orgre.organization_id=o.id
			JOIN employee e ON orgre.user_id=e.id
			WHERE e.username=$1;
		
		`
	} else if bidAuthorType == models.BidAuthorUser {
		query = `SELECT id FROM employee WHERE username=$1`
	}

	if err = db.QueryRow(query, username).Scan(&authorID); err != nil {
		http.Error(w, "Can't find author id", http.StatusNotFound)
		return
	}

	query = `
		SELECT status FROM bid
		WHERE id=$1 AND author_id=$2
	`
	var bidStatus string
	if err := db.QueryRow(query, bidID, authorID).Scan(&bidStatus); err != nil {
		http.Error(w, "AuthorID incorrect or not responsible organization", http.StatusInternalServerError)
		return
	}

	if status == bidStatus {
		http.Error(w, "Status is already up to date", http.StatusConflict)
		return
	}

	var bid models.Bid
	query = `UPDATE bid SET status=$1 WHERE id=$2 RETURNING *`
	if err := db.QueryRow(query, status, bidID).Scan(&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID, &bid.AuthorType, &bid.AuthorID, &bid.Version, &bid.CreatedAt); err != nil {
		http.Error(w, "Failed to update tender status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bid)
}

func EditBidHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID := vars["bidId"]

	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username не указан", http.StatusBadRequest)
		return
	}

	var bidUpdate models.BidUpdate
	if err := json.NewDecoder(r.Body).Decode(&bidUpdate); err != nil {
		http.Error(w, "Неправильно сформированы данные", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	var bidAuthorType models.BidAuthorType
	query := `SELECT author_type FROM bid WHERE id=$1`
	if err := db.QueryRow(query, bidID).Scan(&bidAuthorType); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var authorID uuid.UUID
	if bidAuthorType == models.BidAuthorOrganization {
		query = `
			SELECT o.id FROM organization o
			JOIN organization_responsible orgre ON orgre.organization_id=o.id
			JOIN employee e ON orgre.user_id=e.id
			WHERE e.username=$1;
		
		`
	} else if bidAuthorType == models.BidAuthorUser {
		query = `SELECT id FROM employee WHERE username=$1`
	}

	if err = db.QueryRow(query, username).Scan(&authorID); err != nil {
		http.Error(w, "Can't find author id", http.StatusNotFound)
		return
	}

	query = `
		UPDATE bid SET name=$1, description=$2
		WHERE id=$3 AND author_id=$4
		RETURNING *
	`

	var bid models.Bid
	if err := db.QueryRow(query, bidUpdate.Name, bidUpdate.Description, bidID, authorID).Scan(&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID, &bid.AuthorType, &bid.AuthorID, &bid.Version, &bid.CreatedAt); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bid)
}

func SubmitBidDecisionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID := vars["bidId"]

	decision := r.URL.Query().Get("decision")
	if decision == "" {
		http.Error(w, "decision не указан", http.StatusBadRequest)
		return
	}

	username := vars["username"]
	if username == "" {
		http.Error(w, "username не указан", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	var bidAuthorType models.BidAuthorType
	query := `SELECT author_type FROM bid WHERE id=$1`
	if err := db.QueryRow(query, bidID).Scan(&bidAuthorType); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var authorID uuid.UUID
	if bidAuthorType == models.BidAuthorOrganization {
		query = `
			SELECT o.id FROM organization o
			JOIN organization_responsible orgre ON orgre.organization_id=o.id
			JOIN employee e ON orgre.user_id=e.id
			WHERE e.username=$1;
		
		`
	} else if bidAuthorType == models.BidAuthorUser {
		query = `SELECT id FROM employee WHERE username=$1`
	}

	if err = db.QueryRow(query, username).Scan(&authorID); err != nil {
		http.Error(w, "Can't find author id", http.StatusNotFound)
		return
	}

	query = `
		UPDATE bid SET status='Closed'
		WHERE id=$1 AND author_id=$2
		RETURNING *
	`

	var bid models.Bid
	if err := db.QueryRow(query, decision, bidID, username).Scan(&bid.ID, &bid.TenderID, &bid.Name, &bid.Name, &bid.Description, &bid.CreatedAt); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bid)
}

func SubmitBidFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
}

func RollbackBidHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
}

func GetBidReviewsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
}
