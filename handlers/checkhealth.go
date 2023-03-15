package handlers

import (
	"fmt"
	"net/http"
)

func (app *App) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}
