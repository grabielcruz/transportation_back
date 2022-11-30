package errors_handler

import "log"

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
