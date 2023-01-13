package persons

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/common"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
)

func GetPersons() []Person {
	persons := []Person{}
	rows, err := database.DB.Query("SELECT * FROM persons WHERE id <> $1;", uuid.UUID{})
	if err != nil {
		errors_handler.HandleError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Person
		err := rows.Scan(&p.ID, &p.Name, &p.Document, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			errors_handler.HandleError(err)
		}
		persons = append(persons, p)
	}
	if err != nil {
		errors_handler.HandleError(rows.Err())
	}
	return persons
}

func CreatePerson(fields PersonFields) (Person, error) {
	p := Person{}
	row := database.DB.QueryRow(
		"INSERT INTO persons (name, document) VALUES ($1, $2) RETURNING *;",
		fields.Name, fields.Document)
	err := row.Scan(&p.ID, &p.Name, &p.Document, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return p, errors_handler.MapDBErrors(err)
	}
	return p, nil
}

func GetOnePerson(person_id uuid.UUID) (Person, error) {
	p := Person{}
	if person_id == (uuid.UUID{}) {
		return p, fmt.Errorf(errors_handler.DB001)
	}
	row := database.DB.QueryRow("SELECT * FROM persons WHERE id = $1;", person_id)
	err := row.Scan(&p.ID, &p.Name, &p.Document, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return p, errors_handler.MapDBErrors(err)
	}
	return p, nil
}

func UpdatePerson(person_id uuid.UUID, fields PersonFields) (Person, error) {
	p := Person{}
	if person_id == (uuid.UUID{}) {
		return p, fmt.Errorf(errors_handler.DB001)
	}
	row := database.DB.QueryRow("UPDATE persons SET name = $1, document = $2, updated_at = $3 WHERE id = $4 RETURNING *;",
		fields.Name, fields.Document, time.Now(), person_id)
	err := row.Scan(&p.ID, &p.Name, &p.Document, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return p, errors_handler.MapDBErrors(err)
	}
	return p, nil
}

func DeleteOnePerson(person_id uuid.UUID) (common.ID, error) {
	id := common.ID{}
	if person_id == (uuid.UUID{}) {
		return id, fmt.Errorf(errors_handler.DB001)
	}
	row := database.DB.QueryRow("DELETE FROM persons WHERE id = $1 RETURNING id;", person_id)
	err := row.Scan(&id.ID)
	if err != nil {
		return id, errors_handler.MapDBErrors(err)
	}
	return id, nil
}

func GetPersonsName(person_id uuid.UUID) (string, error) {
	var name string = ""
	if person_id == (uuid.UUID{}) {
		return name, fmt.Errorf(errors_handler.DB001)
	}
	row := database.DB.QueryRow("SELECT name FROM persons WHERE id = $1;", person_id)
	err := row.Scan(&name)
	if err != nil {
		return name, errors_handler.MapDBErrors(err)
	}
	return name, nil
}

func DeleteAllPersons() {
	database.DB.QueryRow("DELETE FROM persons WHERE id <> $1;", uuid.UUID{})
}
