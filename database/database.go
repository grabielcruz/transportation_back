package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/grabielcruz/transportation_back/environment"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func SetupDB() {
	var err error

	myEnv := environment.LoadEnvironment()

	host := myEnv["host"]
	strPort := myEnv["port"]
	user := myEnv["user"]
	password := myEnv["password"]
	DBname := myEnv["DBname"]

	port, _ := strconv.Atoi(strPort)
	errors_handler.CheckError(err)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s DBname=%s sslmode=disable", host, port, user, password, DBname)
	DB, err = sql.Open("postgres", psqlInfo)
	errors_handler.CheckError(err)

	err = DB.Ping()
	errors_handler.CheckError(err)

	log.Println("Database connected")
}

func GetDB() *sql.DB {
	return DB
}

func ResetDB

func CloseConnection() error {
	return DB.Close()
}
