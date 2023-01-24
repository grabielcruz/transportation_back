package person_accounts

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/common"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
)

func GetPersonAccounts(person_id uuid.UUID) ([]PersonAccount, error) {
	accounts := []PersonAccount{}
	rows, err := database.DB.Query("SELECT * FROM person_accounts WHERE person_id = $1;", person_id)
	if err != nil {
		errors_handler.HandleError(err)
	}
	for rows.Next() {
		p := PersonAccount{}
		err := rows.Scan(&p.ID, &p.PersonId, &p.Name, &p.Description, &p.Currency, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return accounts, errors_handler.MapDBErrors(err)
		}
		accounts = append(accounts, p)
	}
	return accounts, nil
}

func CreatePersonAccount(person_id uuid.UUID, fields PersonAccountFields) (PersonAccount, error) {
	pa := PersonAccount{}
	row := database.DB.QueryRow(
		"INSERT INTO person_accounts (person_id, name, description, currency) VALUES ($1, $2, $3, $4) RETURNING *;",
		person_id, fields.Name, fields.Description, fields.Currency)
	err := row.Scan(&pa.ID, &pa.PersonId, &pa.Name, &pa.Description, &pa.Currency, &pa.CreatedAt, &pa.UpdatedAt)
	if err != nil {
		return pa, errors_handler.MapDBErrors(err)
	}
	return pa, nil
}

func GetOnePersonAccount(person_account_id uuid.UUID) (PersonAccount, error) {
	pa := PersonAccount{}
	row := database.DB.QueryRow("SELECT * FROM person_accounts WHERE id = $1;", person_account_id)
	err := row.Scan(&pa.ID, &pa.PersonId, &pa.Name, &pa.Description, &pa.Currency, &pa.CreatedAt, &pa.UpdatedAt)
	if err != nil {
		return pa, fmt.Errorf(errors_handler.PA002)
	}
	return pa, nil
}

// UpdatePersonAccount won't update currency, only name and description
func UpdatePersonAccount(person_account_id uuid.UUID, fields UpdatePersonAccountFields) (PersonAccount, error) {
	pa := PersonAccount{}
	row := database.DB.QueryRow("UPDATE person_accounts SET name = $1, description = $2, updated_at = $3 WHERE id = $4 RETURNING *;",
		fields.Name, fields.Description, time.Now(), person_account_id)
	err := row.Scan(&pa.ID, &pa.PersonId, &pa.Name, &pa.Description, &pa.Currency, &pa.CreatedAt, &pa.UpdatedAt)
	if err != nil {
		return pa, errors_handler.MapDBErrors(err)
	}
	return pa, nil
}

func DeletePersonAccount(person_account_id uuid.UUID) (common.ID, error) {
	id := common.ID{}
	row := database.DB.QueryRow("DELETE FROM person_accounts WHERE id = $1 RETURNING id;", person_account_id)
	err := row.Scan(&id.ID)
	if err != nil {
		return id, errors_handler.MapDBErrors(err)
	}
	return id, nil
}

func DeleteAllPersonAccounts() {
	database.DB.QueryRow("DELETE FROM person_accounts;")
}
