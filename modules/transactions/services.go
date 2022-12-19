package transactions

import (
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
)

func GetTransactions(offset int, limit int) TransationResponse {
	tr := TransationResponse{}
	return tr
}

func CreateTransaction(fields TransactionFields) (Transaction, error) {
	tr := Transaction{}
	balance := 0
	row := database.DB.QueryRow(
		"INSERT INTO transactions (account_id, person_id, date, amount, description, balance) VALUES ($1, $2, $3, $4, $5, %6 RETURNING *;",
		fields.AccountId, fields.PersonId, fields.Date, fields.Amount, fields.Description, balance)
	err := row.Scan(&tr.ID, &tr.AccountId, &tr.PersonId, &tr.Date, &tr.Amount, &tr.Description, &tr.Balance, &tr.CreatedAt, &tr.UpdatedAt)
	errors_handler.CheckError(err)
	return tr, nil
}
