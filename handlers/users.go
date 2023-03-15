package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/carlosarraes/fsback/utils"
	"github.com/go-chi/chi/v5"
)

type User struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Progress  float64 `json:"progress"`
}

func (app *App) GetUsers(w http.ResponseWriter, r *http.Request) {
	data, err := app.Db.Query("SELECT first_name, last_name, progress FROM data.user")
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Error getting users")
		return
	}
	defer data.Close()

	var usersList []User
	for data.Next() {
		var user User
		err := data.Scan(&user.FirstName, &user.LastName, &user.Progress)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, "Error getting users")
			return
		}
		usersList = append(usersList, user)
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(usersList); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Error getting users")
		return
	}
}

func (app *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
	lastName := chi.URLParam(r, "lastName")

	data, err := app.Db.Exec("DELETE FROM data.user WHERE last_name = $1", lastName)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Error deleting user")
		return
	}
	status, _ := data.RowsAffected()
	if status == 0 {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Error deleting user: User not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Error creating user: Invalid request body")
		return
	}

	if user.FirstName == "" || user.LastName == "" || user.Progress == 0 {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Error creating user: First name, last name and progress are required")
		return
	}

	sumCheck, err := app.Db.Query("SELECT sum(progress) FROM data.user")
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Error getting users")
		return
	}
	defer sumCheck.Close()
	var sum float64
	if sumCheck.Next() {
		err := sumCheck.Scan(&sum)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, "Error scanning users")
			return
		}
	}

	if (sum*100)+user.Progress > 100 {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Error creating user: Progress sum is greater than 100")
		return
	}

	data, err := app.Db.Exec("INSERT INTO data.user (first_name, last_name, progress) VALUES ($1, $2, $3)", user.FirstName, user.LastName, user.Progress)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Error creating user")
		return
	}
	status, _ := data.RowsAffected()
	if status == 0 {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	w.WriteHeader(http.StatusCreated)
}
