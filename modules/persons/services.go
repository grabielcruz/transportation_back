package persons

import (
	"time"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/common"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
)

func GetPersons() []Person {
	persons := []Person{}
	rows, err := database.DB.Query("SELECT * FROM persons;")
	errors_handler.CheckError(err)
	defer rows.Close()

	for rows.Next() {
		var p Person
		err := rows.Scan(&p.ID, &p.Name, &p.Document, &p.CreatedAt, &p.UpdatedAt)
		errors_handler.CheckError(err)
		persons = append(persons, p)
	}
	errors_handler.CheckError(rows.Err())
	return persons
}

func CreatePerson(fields PersonFields) Person {
	p := Person{}
	row := database.DB.QueryRow(
		"INSERT INTO persons (name, document) VALUES ($1, $2) RETURNING *;",
		fields.Name, fields.Document)
	err := row.Scan(&p.ID, &p.Name, &p.Document, &p.CreatedAt, &p.UpdatedAt)
	errors_handler.CheckError(err)
	return p
}

func GetOnePerson(person_id uuid.UUID) (Person, error) {
	p := Person{}
	row := database.DB.QueryRow("SELECT * FROM persons WHERE id = $1;", person_id)
	err := row.Scan(&p.ID, &p.Name, &p.Document, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors_handler.CheckEmptyRowError(err) {
			return p, err
		}
		errors_handler.CheckError(err)
	}
	return p, nil
}

func UpdatePerson(person_id uuid.UUID, fields PersonFields) (Person, error) {
	p := Person{}
	row := database.DB.QueryRow("UPDATE persons SET name = $1, document = $2, updated_at = $3 WHERE id = $4 RETURNING *;",
		fields.Name, fields.Document, time.Now(), person_id)
	err := row.Scan(&p.ID, &p.Name, &p.Document, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors_handler.CheckEmptyRowError(err) {
			return p, err
		}
		errors_handler.CheckError(err)
	}
	return p, nil
}

func DeleteOnePerson(person_id uuid.UUID) (common.ID, error) {
	id := common.ID{}
	row := database.DB.QueryRow("DELETE FROM persons WHERE id = $1 RETURNING id;", person_id)
	err := row.Scan(&id.ID)
	if err != nil {
		if errors_handler.CheckEmptyRowError(err) {
			return id, err
		}
		errors_handler.CheckError(err)
	}
	return id, nil
}

func deleteAllPersons() {
	database.DB.QueryRow("DELETE FROM persons;")
}
