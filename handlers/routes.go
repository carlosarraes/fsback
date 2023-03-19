package handlers

import (
	"database/sql"
	"net/http"

	"github.com/carlosarraes/fsback/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	Db *sql.DB
}

func (a *App) Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(utils.Cors())

	mux.Get("/", a.CheckHealth)
	mux.Get("/users", a.GetUsers)
	mux.Post("/user", a.CreateUser)
	mux.Delete("/user/{lastName}", a.DeleteUser)

	return mux
}
