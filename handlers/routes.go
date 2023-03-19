package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/carlosarraes/fsback/db"
	"github.com/carlosarraes/fsback/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	DSN string
	DB  db.PostgresConn
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

func (a *App) Connect() (*sql.DB, error) {
	connection, err := db.OpenDB(a.DSN)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to db!")
	return connection, nil
}
