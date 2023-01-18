package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/grabielcruz/transportation_back/environment"
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

	port, err := strconv.Atoi(strPort)
	if err != nil {
		log.Fatal(err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	fmt.Println(psqlInfo)
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connected")
}

func GetDB() *sql.DB {
	return DB
}

func CreateTables(sqlPath string) {
	dat, err := os.ReadFile(sqlPath)
	if err != nil {
		log.Fatal(err)
	}
	sqlStr := string(dat)
	_, err = DB.Exec(sqlStr)
	if err != nil {
		log.Fatal(err)
	}
}

func CloseConnection() error {
	return DB.Close()
}
