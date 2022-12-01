package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/grabielcruz/transportation_back/environment"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func SetupDB(mode string) {
	var err error

	myEnv := environment.LoadEnvironment(mode)

	host := myEnv["host"]
	strPort := myEnv["port"]
	user := myEnv["user"]
	password := myEnv["password"]
	dbname := myEnv["dbname"]

	port, _ := strconv.Atoi(strPort)
	errors_handler.CheckError(err)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	fmt.Println(psqlInfo)
	DB, err = sql.Open("postgres", psqlInfo)
	errors_handler.CheckError(err)

	err = DB.Ping()
	errors_handler.CheckError(err)

	log.Println("Database connected")
}

func GetDB() *sql.DB {
	return DB
}

func CreateTables(sqlPath string) {
	dat, err := os.ReadFile(sqlPath)
	errors_handler.CheckError(err)
	sqlStr := string(dat)
	_, err = DB.Exec(sqlStr)
	errors_handler.CheckError(err)
}

func CloseConnection() error {
	return DB.Close()
}
