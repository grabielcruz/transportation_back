package environment

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnvironment
func LoadEnvironment(envPath string) map[string]string {
	var myEnv map[string]string
	myEnv, err := godotenv.Read(envPath)
	if err != nil {
		log.Fatal(err)
	}
	return myEnv
}
