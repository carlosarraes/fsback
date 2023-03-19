package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/carlosarraes/fsback/db"
	"github.com/carlosarraes/fsback/utils"
	"github.com/go-chi/chi/v5"
)

func (app *App) GetUsers(w http.ResponseWriter, r *http.Request) {
	usersList, err := app.DB.GetUsers()
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, "Error getting users")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(usersList); err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, "Error getting users")
		return
	}
}

func (app *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
	lastName := chi.URLParam(r, "lastName")

	err := app.DB.DeleteUser(lastName)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, "Error deleting user")
		return
	}

	user := fmt.Sprintf("User %s deleted", lastName)
	utils.WriteResponse(w, http.StatusOK, user)
}

func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user db.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, "Error creating user: Invalid request body")
		return
	}

	if user.FirstName == "" || user.LastName == "" || user.Progress == 0 {
		utils.WriteResponse(w, http.StatusBadRequest, "Error creating user: First name, last name and progress are required")
		return
	}

	sum, err := app.DB.SumCheck()
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	if sum+user.Progress > 1 {
		utils.WriteResponse(w, http.StatusBadRequest, "Error creating user: Progress sum cannot exceed 100")
		return
	}

	if err = app.DB.CreateUser(user); err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	utils.WriteResponse(w, http.StatusCreated, "User created")
}
