package main

import (
	"path/filepath"

	"github.com/grabielcruz/transportation_back/database"
)

func main() {
	envPath := filepath.Clean(".env")
	database.SetupDB(envPath)
}
