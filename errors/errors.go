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
	case CU001:
		return "CU001"
	case CU002:
		return "CU002"
	case CU003:
		return "CU003"
	case CU004:
		return "CU004"
	case PE001:
		return "PE001"
	case TR001:
		return "TR001"
	case TR002:
		return "TR002"
	case TR003:
		return "TR003"
	case TR004:
		return "TR004"
	case TR005:
		return "TR005"
	case TR006:
		return "TR006"
	case TR007:
		return "TR007"
	case TR008:
		return "TR008"
	case TR009:
		return "TR009"
	case TR010:
		return "TR010"
	case TR011:
		return "TR011"
	default:
		return "SE001"
	}
}
