package errors_handler

import (
	"fmt"
	"log"
	"os"
)

// module level
const ErrorsLogPath = "../../errors_log.txt"

func HandleError(err error) {
	WriteErrorToFile(ErrorsLogPath, err.Error())
}

func WriteErrorToFile(filePath string, msg string) {
	f, err := os.OpenFile(filePath, os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.WriteString(msg + "\n")
	if err != nil {
		log.Fatal(err)
	}
}

func ResetFile(path string) {
	err := os.WriteFile(path, []byte(""), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func MapDBErrors(err error) error {
	switch err.Error() {
	// Database
	case "sql: no rows in result set":
		return fmt.Errorf(DB001)

	// Currencies
	case "pq: duplicate key value violates unique constraint \"currencies_pkey\"":
		return fmt.Errorf(CU003)
	case "pq: update or delete on table \"currencies\" violates foreign key constraint \"money_accounts_currency_fkey\" on table \"money_accounts\"":
		return fmt.Errorf(CU004)

	// Persons
	case "pq: duplicate key value violates unique constraint \"persons_document_key\"":
		return fmt.Errorf(PE001)

	// Bills
	case "pq: insert or update on table \"pending_bills\" violates foreign key constraint \"pending_bills_currency_fkey\"":
		return fmt.Errorf(CU005)
	case "pq: insert or update on table \"closed_bills\" violates foreign key constraint \"pending_bills_currency_fkey\"":
		return fmt.Errorf(CU005)
	case "pq: update or delete on table \"pending_bills\" violates foreign key constraint \"fk_transactions_pending_bills\" on table \"transactions\"":
		return fmt.Errorf(BL003)
	}
	return err
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
	case DB007:
		return "DB007"

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
	case PE002:
		return "PE002"

	// transactions
	case TR001:
		return "TR001"
	case TR002:
		return "TR002"
	case TR003:
		return "TR003"
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
	case TR012:
		return "TR012"
	case TR013:
		return "TR013"
	case TR014:
		return "TR014"

	// bills
	case BL001:
		return "BL001"
	case BL002:
		return "BL002"
	case BL003:
		return "BL003"

	// person accounts
	case PA001:
		return "PA001"
	case PA002:
		return "PA002"

	//default
	default:
		return "SE001"
	}
}
