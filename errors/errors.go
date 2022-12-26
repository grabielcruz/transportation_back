package errors_handler

import "log"

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CheckEmptyRowError(err error) bool {
	return err.Error() == DB001
}

func MapServiceError(error_msg string) string {
	switch error_msg {
	case DB001:
		return "DB001"
	case TR001:
		return "TR001"
	case TR002:
		return "TR002"
	case TR003:
		return "TR003"
	case TR004:
		return "TR004"
	default:
		return "SE001"
	}
}
