package bills

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/modules/persons"
	"github.com/lib/pq"
)

// GetPendingBills returns the pending bills paginated, filtered by person, wether it is to be paid, it is to be charged
// limit and offset are for pagination porpuses
func GetPendingBills(person_id uuid.UUID, to_pay bool, to_charge bool, limit int, offset int) (BillResponse, error) {
	return getBills("pending_bills", person_id, to_pay, to_charge, limit, offset)
}

func CreatePendingBill(fields BillFields) (Bill, error) {
	bill := Bill{}
	if fields.Amount == float64(0) {
		return bill, fmt.Errorf(errors_handler.BL002)
	}
	row := database.DB.QueryRow("INSERT INTO pending_bills (person_id, date, description, currency, amount, pending) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;", fields.PersonId, fields.Date, fields.Description, fields.Currency, fields.Amount, fields.Amount)
	err := row.Scan(&bill.ID, &bill.PersonId, &bill.Date, &bill.Description, &bill.Currency, &bill.Amount, &bill.Pending, &bill.CreatedAt, &bill.UpdatedAt)
	if err != nil {
		// person with given uuid does not exists
		if err.(*pq.Error).Message == "insert or update on table \"pending_bills\" violates foreign key constraint \"pending_bills_person_id_fkey\"" {
			return bill, fmt.Errorf(errors_handler.BL003)
		}
		// currency not registered
		if err.(*pq.Error).Message == "insert or update on table \"pending_bills\" violates foreign key constraint \"pending_bills_currency_fkey\"" {
			return bill, fmt.Errorf(errors_handler.BL004)
		}
		return bill, fmt.Errorf(errors_handler.DB007)
	}

	bill.PersonName, _ = persons.GetPersonsName(bill.PersonId)

	return bill, nil
}

// GetOneBill

func GetClosedBills() {

}

func getBills(table_name string, person_id uuid.UUID, to_pay bool, to_charge bool, limit int, offset int) (BillResponse, error) {
	billResponse := BillResponse{}
	filters := []string{}

	// can't have to_pay and to_charge on false at the same time
	if !to_pay && !to_charge {
		return billResponse, fmt.Errorf(errors_handler.BL001)
	}

	// check if person_id is not zero uuid
	if person_id.String() != (uuid.UUID{}).String() {
		// should be safe, it is an uuid
		filters = append(filters, fmt.Sprintf("person_id = '%v'", person_id.String()))
	}

	// only one can happen at a time
	if !to_pay {
		filters = append(filters, "amount > 0")
	}

	if !to_charge {
		filters = append(filters, "amount < 0")
	}

	searchString := strings.Join(filters, " AND ")
	if len(searchString) > 0 {
		searchString = "WHERE " + searchString
	}

	tx, err := database.DB.Begin()
	if err != nil {
		tx.Rollback()
		return billResponse, fmt.Errorf(errors_handler.DB002)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %v %v;", table_name, searchString)
	row := tx.QueryRow(countQuery)
	err = row.Scan(&billResponse.Count)
	if err != nil {
		tx.Rollback()
		return billResponse, fmt.Errorf(errors_handler.DB004)
	}

	recordsQuery := fmt.Sprintf("SELECT * FROM %v %v ORDER BY created_at DESC LIMIT $1 OFFSET $2;", table_name, searchString)
	rows, err := tx.Query(recordsQuery, limit, offset)
	if err != nil {
		tx.Rollback()
		return billResponse, fmt.Errorf(errors_handler.DB005)
	}

	for rows.Next() {
		b := Bill{}
		err = rows.Scan(&b.ID, &b.PersonId, &b.Date, &b.Description, &b.Currency, &b.Amount, &b.Pending, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			tx.Rollback()
			return billResponse, fmt.Errorf(errors_handler.DB006)
		}
		b.PersonName, _ = persons.GetPersonsName(person_id)
		billResponse.Bills = append(billResponse.Bills, b)
	}

	billResponse.Limit = limit
	billResponse.Offset = offset
	billResponse.FilterPersonId = person_id

	err = tx.Commit()
	if err != nil {
		return billResponse, fmt.Errorf(errors_handler.DB003)
	}
	return billResponse, nil
}

func emptyBills() {
	database.DB.QueryRow("DELETE FROM pending_bills;")
	database.DB.QueryRow("DELETE FROM closed_bills;")
}
