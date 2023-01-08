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
	// database
	case DB001:
		return "DB001"
	case DB002:
		return "DB002"
	case DB003:
		return "DB003"
	case DB004:
		return "DB004"
	case DB005:
		return "DB005"
	case DB006:
		return "DB006"
	case DB007:
		return "DB007"
	case DB008:
		return "DB008"

	// currencies
	case CU001:
		return "CU001"
	case CU002:
		return "CU002"
	case CU003:
		return "CU003"
	case CU004:
		return "CU004"

	// persons
	case PE001:
		return "PE001"

	// transactions
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
	case TR011:
		return "TR011"
	case TR012:
		return "TR012"

	// bills
	case BL001:
		return "BL001"
	case BL002:
		return "BL002"
	case BL003:
		return "BL003"
	case BL004:
		return "BL004"

	//default
	default:
		return "SE001"
	}
}
