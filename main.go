package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/carlosarraes/fsback/db"
	"github.com/carlosarraes/fsback/handlers"
	"github.com/joho/godotenv"
)

func main() {
	port := 8080
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	app := handlers.App{}

	flag.StringVar(&app.DSN, "dsn", os.Getenv("DB_URI"), "Database connection string")
	flag.Parse()

	conn, err := app.Connect()
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}
	defer conn.Close()

	app.DB = db.PostgresConn{DB: conn}

	r := app.Routes()

	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
