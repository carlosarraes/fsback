package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting users: %v", err)
	}
	defer data.Close()

	var usersList []User
	for data.Next() {
		var user User
		err := data.Scan(&user.FirstName, &user.LastName, &user.Progress)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error scanning users: %v", err)
		}
		usersList = append(usersList, user)
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(usersList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error encoding users: %v", err)
	}
}

func (app *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
	lastName := chi.URLParam(r, "lastName")

	data, err := app.Db.Exec("DELETE FROM data.user WHERE last_name = $1", lastName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting user: %v", err)
		return
	}
	status, _ := data.RowsAffected()
	if status == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
