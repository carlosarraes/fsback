package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/carlosarraes/fsback/db"
	"github.com/carlosarraes/fsback/handlers"
)

func main() {
	port := 8080

	db, err := db.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	app := handlers.App{Db: db}

	r := app.Routes()

	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
