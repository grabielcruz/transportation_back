package environment

import (
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/joho/godotenv"
)

func LoadEnvironment() map[string]string {
	var myEnv map[string]string
	myEnv, err := godotenv.Read("../.env")
	errors_handler.CheckError(err)
	return myEnv
}
