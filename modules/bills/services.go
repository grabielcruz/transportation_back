package bills

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/common"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/modules/persons"
)

// GetPendingBills returns the pending bills paginated, filtered by person, wether it is to be paid, it is to be charged
// limit and offset are for pagination porpuses
func GetPendingBills(person_id uuid.UUID, to_pay bool, to_charge bool, limit int, offset int) (BillResponse, error) {
	billResponse := BillResponse{}
	filters := []string{}
	// to exclude zero bill
	filters = append(filters, "id <> $1")

	// can't have to_pay and to_charge on false at the same time
	if !to_pay && !to_charge {
		return billResponse, fmt.Errorf(errors_handler.BL001)
	}

	// check if person_id is not zero uuid
	if person_id.String() != (uuid.UUID{}).String() {
		// should be safe, it is an uuid
		filters = append(filters, fmt.Sprintf("person_id = '%v'", person_id.String()))
	}

	// only one of these can happen at a time
	if !to_pay {
		filters = append(filters, "amount > 0")
	}

	if !to_charge {
		filters = append(filters, "amount < 0")
	}
	//

	searchString := strings.Join(filters, " AND ")
	if len(searchString) > 0 {
		searchString = "WHERE " + searchString
	}

	tx, err := database.DB.Begin()
	if err != nil {
		tx.Rollback()
		return billResponse, fmt.Errorf(errors_handler.DB002)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM pending_bills %v;", searchString)
	row := tx.QueryRow(countQuery, uuid.UUID{})
	err = row.Scan(&billResponse.Count)
	if err != nil {
		tx.Rollback()
		return billResponse, fmt.Errorf(errors_handler.DB004)
	}

	recordsQuery := fmt.Sprintf("SELECT * FROM pending_bills %v ORDER BY created_at DESC LIMIT $2 OFFSET $3;", searchString)
	rows, err := tx.Query(recordsQuery, uuid.UUID{}, limit, offset)
	if err != nil {
		tx.Rollback()
		return billResponse, fmt.Errorf(errors_handler.DB005)
	}

	for rows.Next() {
		b := Bill{}
		err = rows.Scan(&b.ID, &b.PersonId, &b.Date, &b.Description, &b.Status, &b.Currency, &b.Amount, &b.ParentTransactionId, &b.ParentBillCrossId, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			tx.Rollback()
			return billResponse, fmt.Errorf(errors_handler.DB005)
		}
		b.PersonName, err = persons.GetPersonsName(person_id)
		if err != nil {
			errors_handler.HandleError(err)
		}
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

func CreatePendingBill(fields BillFields) (Bill, error) {
	bill := Bill{}
	if fields.Amount == float64(0) {
		return bill, fmt.Errorf(errors_handler.BL002)
	}
	row := database.DB.QueryRow("INSERT INTO pending_bills (person_id, date, description, currency, amount, parent_transaction_id, parent_bill_cross_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;", fields.PersonId, fields.Date, fields.Description, fields.Currency, fields.Amount, uuid.UUID{}, uuid.UUID{})
	err := row.Scan(&bill.ID, &bill.PersonId, &bill.Date, &bill.Description, &bill.Status, &bill.Currency, &bill.Amount, &bill.ParentTransactionId, &bill.ParentBillCrossId, &bill.CreatedAt, &bill.UpdatedAt)
	if err != nil {
		return bill, errors_handler.MapDBErrors(err)
	}

	bill.PersonName, err = persons.GetPersonsName(bill.PersonId)
	if err != nil {
		errors_handler.HandleError(err)
	}

	return bill, nil
}

func GetOneBill(bill_id uuid.UUID) (Bill, error) {
	b := Bill{}
	if bill_id == (uuid.UUID{}) {
		return b, fmt.Errorf(errors_handler.DB001)
	}
	row := database.DB.QueryRow("SELECT * FROM pending_bills WHERE id = $1;", bill_id)
	err := row.Scan(&b.ID, &b.PersonId, &b.Date, &b.Description, &b.Status, &b.Currency, &b.Amount, &b.ParentTransactionId, &b.ParentBillCrossId, &b.CreatedAt, &b.UpdatedAt)

	// not found in pending_bills, look for it on closed bills
	if err != nil {
		row = database.DB.QueryRow("SELECT * FROM closed_bills WHERE id = $1;", bill_id)
		err = row.Scan(&b.ID, &b.PersonId, &b.Date, &b.Description, &b.Status, &b.Currency, &b.Amount, &b.ParentTransactionId, &b.ParentBillCrossId, &b.TransactionId, &b.BillCrossId, &b.RevertTransactionId, &b.PostNotes, &b.CreatedAt, &b.UpdatedAt)
		// bill not found anywhere
		if err != nil {
			return b, fmt.Errorf(errors_handler.DB001)
		}
	}
	b.PersonName, err = persons.GetPersonsName(b.PersonId)
	if err != nil {
		errors_handler.HandleError(err)
	}
	return b, nil
}

func UpdatePendingBill(bill_id uuid.UUID, fields BillFields) (Bill, error) {
	b := Bill{}
	if fields.PersonId == (uuid.UUID{}) {
		return b, fmt.Errorf(errors_handler.PE002)
	}
	if bill_id == (uuid.UUID{}) {
		return b, fmt.Errorf(errors_handler.DB001)
	}
	row := database.DB.QueryRow("UPDATE pending_bills SET person_id = $1, date = $2, description = $3, currency = $4, amount = $5 WHERE id = $6 RETURNING *;", fields.PersonId, fields.Date, fields.Description, fields.Currency, fields.Amount, bill_id)
	err := row.Scan(&b.ID, &b.PersonId, &b.Date, &b.Description, &b.Status, &b.Currency, &b.Amount, &b.ParentTransactionId, &b.ParentBillCrossId, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return b, errors_handler.MapDBErrors(err)
	}
	b.PersonName, err = persons.GetPersonsName(b.PersonId)
	if err != nil {
		errors_handler.HandleError(err)
	}
	return b, nil
}

func DeleteBill(bill_id uuid.UUID) (common.ID, error) {
	id := common.ID{}
	if bill_id == (uuid.UUID{}) {
		return id, fmt.Errorf(errors_handler.DB001)
	}
	row := database.DB.QueryRow("DELETE FROM pending_bills WHERE id = $1 RETURNING id;", bill_id)
	err := row.Scan(&id.ID)
	if err != nil {
		return id, errors_handler.MapDBErrors(err)
	}
	return id, nil
}

func createClosedBill(fields BillFields) (Bill, error) {
	bill := Bill{}
	if fields.Amount == float64(0) {
		return bill, fmt.Errorf(errors_handler.BL002)
	}
	randomUUID, _ := uuid.NewRandom()
	row := database.DB.QueryRow("INSERT INTO closed_bills (id, person_id, date, description, currency, amount, parent_transaction_id, parent_bill_cross_id, transaction_id, bill_cross_id, post_notes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING *;", randomUUID, fields.PersonId, fields.Date, fields.Description, fields.Currency, fields.Amount, uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, "")
	err := row.Scan(&bill.ID, &bill.PersonId, &bill.Date, &bill.Description, &bill.Status, &bill.Currency, &bill.Amount, &bill.ParentTransactionId, &bill.ParentBillCrossId, &bill.TransactionId, &bill.BillCrossId, &bill.RevertTransactionId, &bill.PostNotes, &bill.CreatedAt, &bill.UpdatedAt)
	if err != nil {
		return bill, errors_handler.MapDBErrors(err)
	}

	bill.PersonName, err = persons.GetPersonsName(bill.PersonId)
	if err != nil {
		errors_handler.HandleError(err)
	}

	return bill, nil
}

func EmptyBills() {
	database.DB.QueryRow("DELETE FROM pending_bills WHERE id <> $1;", uuid.UUID{})
	database.DB.QueryRow("DELETE FROM closed_bills WHERE id <> $1;", uuid.UUID{})
}
