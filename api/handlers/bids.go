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

func CreateBidHandler(w http.ResponseWriter, r *http.Request) {
	var bid models.Bid
	if err := json.NewDecoder(r.Body).Decode(&bid); err != nil {
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

	// var authorID uuid.UUID
	// query := `SELECT id FROM employee WHERE username=$1;`
	// if err := db.QueryRow(query, tenderCreate.CreatorUsername).Scan(&authorID); err != nil {
	// 	log.Error(err.Error())
	// 	http.Error(w, fmt.Sprintf("No user with username: {%s}", tenderCreate.CreatorUsername), http.StatusNotFound)
	// 	return
	// }

	// var exists bool
	// query = `SELECT EXISTS(SELECT 1 FROM employee WHERE id=$1);`
	// err = db.QueryRow(query, tenderCreate.OrganizationID).Scan(&exists)
	// if err != nil || !exists {
	// 	log.Error(err.Error())
	// 	http.Error(w, fmt.Sprintf("No organization with id: {%s}", tenderCreate.OrganizationID), http.StatusNotAcceptable)
	// 	return
	// }

	// query = `SELECT EXISTS(SELECT 1 FROM organization_responsible WHERE organization_id=$1 and user_id=$2);`
	// if err := db.QueryRow(query, tenderCreate.OrganizationId).Scan(&organizationID); err != nil {
	// 	log.Error(err.Error())
	// 	http.Error(w, fmt.Sprintf("No user with username: {%s}", tenderCreate.CreatorUsername), http.StatusNotFound)
	// 	return
	// }

	query := `
		INSERT INTO bid (name, description, status, tender_id, organization_id, creatorUsername)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at;
	`

	if err := db.QueryRow(query, bid.Name, bid.Description, bid.Status, bid.TenderID, bid.AuthorType, bid.AuthorID, bid.Version, time.Now()).Scan(&bid.ID, &bid.CreatedAt); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	log.Infof("Created new bid with ID{%s}", bid.ID)

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

	var bids []models.Bid
	query := `
		SELECT b.* FROM bid b
		JOIN organization o ON b.tender_id IN (SELECT id FROM tender WHERE organization_id = o.id)
		JOIN employee e ON o.id = e.organization_id
		WHERE e.username = $1
		ORDER BY b.name LIMIT $2 OFFSET $3
	`

	rows, err := db.Query(query, username, limit, offset)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Failed to get bids", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var bid models.Bid
		if err := rows.Scan(&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID, &bid.AuthorType, &bid.AuthorID, &bid.Version, &bid.CreatedAt); err != nil {
			log.Error(err.Error())
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

	username := r.URL.Query().Get("username")
	if username == "" {
		log.Error("username не указан")
		http.Error(w, "username не указан", http.StatusBadRequest)
		return
	}

	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	// Проверка авторизации и прав доступа
	// ...

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	query := `
		SELECT * FROM bids
		WHERE tender_id=$1 AND name=$2
		ORDER BY id
		LIMIT $3
	`

	rows, err := db.Query(query, tenderID, username, limit, offset)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bids []models.Bid
	for rows.Next() {
		var bid models.Bid
		if err := rows.Scan(&bid.ID, &bid.TenderID, &bid.Name, &bid.CreatedAt); err != nil {
			log.Error(err.Error())
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
		log.Error("username не указан")
		http.Error(w, "username не указан", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	query := `
		SELECT status FROM bids
		WHERE id=$1 AND username=$2
	`

	var status string
	if err := db.QueryRow(query, bidID, username).Scan(&status); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

func UpdateBidStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	bidID := vars["bidId"]
	name := r.URL.Query().Get("username")

	if name == "" {
		log.Error("username не указан")
		http.Error(w, "username не указан", http.StatusBadRequest)
		return
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		log.Error("status не указан")
		http.Error(w, "status не указан", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	query := `
		UPDATE bids SET status=$1
		WHERE id=$2 AND username=$3
		RETURNING *
	`

	var bid models.Bid
	if err := db.QueryRow(query, status, bidID, name).Scan(&bid.ID, &bid.TenderID, &bid.Name, &bid.Status, &bid.CreatedAt); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
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
		log.Error("Name не указан")
		http.Error(w, "username не указан", http.StatusBadRequest)
		return
	}

	var updatedBid models.Bid
	if err := json.NewDecoder(r.Body).Decode(&updatedBid); err != nil {
		log.Error(err.Error())
		http.Error(w, "Неправильно сформированы данные", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	query := `
		UPDATE bids SET
			name = COALESCE($1, name),
			description = COALESCE($2, description)
		WHERE id=$3 AND username=$4
		RETURNING *
	`

	var bid models.Bid
	if err := db.QueryRow(query, updatedBid.Name, updatedBid.Description, bidID, username).Scan(&bid.ID, &bid.TenderID, &bid.Name, &bid.Name, &bid.Description, &bid.CreatedAt); err != nil {
		log.Error(err.Error())
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
		log.Error("decision не указан")
		http.Error(w, "decision не указан", http.StatusBadRequest)
		return
	}

	username := vars["username"]
	if username == "" {
		log.Error("username не указан")
		http.Error(w, "username не указан", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	query := `
		UPDATE bids SET
			decision = $1
		WHERE id=$2 AND username=$3
		RETURNING *
	`

	var bid models.Bid
	if err := db.QueryRow(query, decision, bidID, username).Scan(&bid.ID, &bid.TenderID, &bid.Name, &bid.Name, &bid.Description, &bid.CreatedAt); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bid)
}

func SubmitBidFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID := vars["bidId"]

	bidFeedback := r.URL.Query().Get("bidFeedback")
	if bidFeedback == "" {
		log.Error("bidFeedback не указан")
		http.Error(w, "bidFeedback не указан", http.StatusBadRequest)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		log.Error("username не указан")
		http.Error(w, "username не указан", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	query := `
		UPDATE bids SET
			feedback = $1
		WHERE id=$2 AND username=$3
		RETURNING *
	`

	var bid models.Bid
	if err := db.QueryRow(query, bidFeedback, bidID, username).Scan(&bid.ID, &bid.TenderID, &bid.AuthorID, &bid.Name, &bid.Description, &bid.Status, &bid.CreatedAt); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bid)
}

func RollbackBidHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID := vars["bidId"]
	version := vars["version"]

	username := r.URL.Query().Get("username")
	if username == "" {
		log.Error("username не указан")
		http.Error(w, "username не указан", http.StatusBadRequest)
		return
	}

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Неправильный формат версии", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	// Получить данные предложения по указанной версии
	query := `
		SELECT * FROM bids
		WHERE id=$1 AND version=$2
	`

	var bid models.Bid
	if err := db.QueryRow(query, bidID, versionInt).Scan(&bid.ID, &bid.TenderID, &bid.AuthorID, &bid.Name, &bid.Description, &bid.Status, &bid.Version, &bid.CreatedAt); err != nil {
		log.Error(err.Error())
		http.Error(w, "Предложение или версия не найдены", http.StatusNotFound)
		return
	}

	// Обновить данные предложения
	query = `
		UPDATE bids SET
			name = $1,
			description = $2,
			status = $3,
			version = version + 1
		WHERE id=$5 AND username=$6
		RETURNING *
	`

	if err := db.QueryRow(query, bid.Name, bid.Description, bid.Status, bidID, username).Scan(&bid.ID, &bid.TenderID, &bid.AuthorID, &bid.Name, &bid.Description, &bid.Status, &bid.Version, &bid.CreatedAt); err != nil {
		log.Error(err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bid)
}

func GetBidReviewsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderID := vars["tenderId"]

	authorUsername := r.URL.Query().Get("authorUsername")
	if authorUsername == "" {
		log.Error("authorUsername не указан")
		http.Error(w, "authorUsername не указан", http.StatusBadRequest)
		return
	}

	requesterUsername := r.URL.Query().Get("requesterUsername")
	if requesterUsername == "" {
		log.Error("requesterUsername не указан")
		http.Error(w, "requesterUsername не указан", http.StatusBadRequest)
		return
	}

	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	db, err := database.Connect()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	query := `
		SELECT * FROM bid_reviews
		WHERE tender_id=$1 AND author_username=$2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	var reviews []models.Review
	if limit != "" && offset != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "Неправильный формат limit", http.StatusBadRequest)
			return
		}

		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "Неправильный формат offset", http.StatusBadRequest)
			return
		}

		rows, err := db.Query(query, tenderID, authorUsername, limitInt, offsetInt)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var review models.Review
			if err := rows.Scan(&review.ID, &review.Description, &review.CreatedAt); err != nil {
				log.Error(err.Error())
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
			reviews = append(reviews, review)
		}
	} else {
		rows, err := db.Query(query, tenderID, authorUsername)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var review models.Review
			if err := rows.Scan(&review.ID, &review.Description, &review.CreatedAt); err != nil {
				log.Error(err.Error())
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
			reviews = append(reviews, review)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reviews)
}
