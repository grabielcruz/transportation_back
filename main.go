package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/grabielcruz/transportation_back/database"
	"github.com/grabielcruz/transportation_back/routes"
)

func main() {
	envPath := filepath.Clean(".env")
	database.SetupDB(envPath)
	sqlPath := filepath.Clean("database/database.sql")
	database.CreateTables(sqlPath)
	defer database.CloseConnection()

	router := routes.SetupAndGetRoutes()

	log.Fatal(http.ListenAndServe(":8080", router))
}
