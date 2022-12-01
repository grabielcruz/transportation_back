package environment

import (
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/joho/godotenv"
)

// LoadEnvironment
func LoadEnvironment(envPath string) map[string]string {
	var myEnv map[string]string
	myEnv, err := godotenv.Read(envPath)
	errors_handler.CheckError(err)
	return myEnv
}
