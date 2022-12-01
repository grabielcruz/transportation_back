package errors_handler

import "log"

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CheckEmptyRowError(err error) bool {
	return err.Error() == "sql: no rows in result set"
}
