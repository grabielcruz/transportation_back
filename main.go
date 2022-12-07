package main

import (
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

	r := routes.SetupAndGetRoutes()

	r.Run()
}
