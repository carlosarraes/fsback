package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/carlosarraes/fsback/db"
	"github.com/carlosarraes/fsback/handlers"
	"github.com/carlosarraes/fsback/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	port := 8080

	db, err := db.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	app := &handlers.App{Db: db}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(utils.Cors())

	r.Get("/", app.CheckHealth)
	r.Get("/users", app.GetUsers)
	r.Delete("/user/{lastName}", app.DeleteUser)

	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
