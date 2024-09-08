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

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `
		INSERT INTO employee(username, first_name, last_name)
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
	if err = db.QueryRow(sqlStatement, user.Username, user.FirstName, user.LastName).Scan(&id, &createdAt, &updatedAt); err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("Created new user with ID{%d}", id)

	user.ID = id
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	sqlStatement := `
		SELECT * FROM employee;
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

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var user models.User
	sqlStatement := `
		SELECT * FROM employee WHERE id=$1
	`

	if err := db.QueryRow(sqlStatement, userID).Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt); err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func UpdateUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	updatedUser.UpdatedAt = time.Now()
	sqlStatement := `
		UPDATE employee SET username=$1, first_name=$2, last_name=$3, updated_at=$4 WHERE id=$5
	`

	if _, err := db.Exec(sqlStatement, updatedUser.Username, updatedUser.FirstName, updatedUser.LastName, updatedUser.UpdatedAt, userID); err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	log.Infof("Updated user{%d} info", userID)

	w.WriteHeader(http.StatusNoContent)
}

func DeleteUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStatement := `
		DELETE FROM employee WHERE id=$1
	`

	if _, err := db.Exec(sqlStatement, userID); err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	log.Infof("Deleted user{%d}", userID)

	w.WriteHeader(http.StatusNoContent)
}
