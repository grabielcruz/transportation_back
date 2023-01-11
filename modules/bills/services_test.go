package bills

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/modules/config"
	"github.com/grabielcruz/transportation_back/modules/persons"
	"github.com/grabielcruz/transportation_back/utility"
	"github.com/stretchr/testify/assert"
)

func TestBillServices(t *testing.T) {
	envPath := filepath.Clean("../../.env_test")
	sqlPath := filepath.Clean("../../database/database.sql")
	database.SetupDB(envPath)
	database.CreateTables(sqlPath)
	defer database.CloseConnection()
	person1, err := persons.CreatePerson(persons.GeneratePersonFields())
	assert.Nil(t, err)
	person2, err := persons.CreatePerson(persons.GeneratePersonFields())
	assert.Nil(t, err)

	t.Run("Get all bills response with zero bills", func(t *testing.T) {
		billResponse, err := GetPendingBills(uuid.UUID{}, true, true, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 0)
		assert.Equal(t, billResponse.Count, 0)
		assert.Equal(t, billResponse.FilterPersonId, uuid.UUID{})
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
	})

	t.Run("Create one pending bill", func(t *testing.T) {
		billFields := GenerateBillFields(person1.ID)
		assert.Nil(t, err)
		newBill, err := CreatePendingBill(billFields, uuid.UUID{}, uuid.UUID{})
		assert.Nil(t, err)
		assert.Equal(t, billFields.PersonId, newBill.PersonId)
		assert.Equal(t, person1.Name, newBill.PersonName)
		assert.Equal(t, billFields.Currency, newBill.Currency)
		assert.Equal(t, billFields.Date.Format("2006-01-02"), newBill.Date.Format("2006-01-02"))
		assert.Equal(t, billFields.Description, newBill.Description)
		assert.Equal(t, billFields.Amount, newBill.Amount)
	})

	emptyBills()

	t.Run("Create 4 bills, 2 for person1, 2 for person2, negative and positive balance and get them filtered", func(t *testing.T) {
		// person1
		billFields := GenerateBillFields(person1.ID)
		billFields.Amount = 55.55
		_, err := CreatePendingBill(billFields, uuid.UUID{}, uuid.UUID{})
		assert.Nil(t, err)

		billFields = GenerateBillFields(person1.ID)
		billFields.Amount = -55.55
		_, err = CreatePendingBill(billFields, uuid.UUID{}, uuid.UUID{})
		assert.Nil(t, err)

		// person2
		billFields = GenerateBillFields(person2.ID)
		billFields.Amount = 77.77
		_, err = CreatePendingBill(billFields, uuid.UUID{}, uuid.UUID{})
		assert.Nil(t, err)

		billFields = GenerateBillFields(person2.ID)
		billFields.Amount = -77.77
		_, err = CreatePendingBill(billFields, uuid.UUID{}, uuid.UUID{})
		assert.Nil(t, err)

		// all of them
		billResponse, err := GetPendingBills(uuid.UUID{}, true, true, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 4)
		assert.Equal(t, billResponse.Count, 4)
		assert.Equal(t, billResponse.FilterPersonId, uuid.UUID{})
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person2.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(-77.77), billResponse.Bills[0].Amount)
		assert.Equal(t, person2.ID, billResponse.Bills[1].PersonId)
		assert.Equal(t, float64(77.77), billResponse.Bills[1].Amount)
		assert.Equal(t, person1.ID, billResponse.Bills[2].PersonId)
		assert.Equal(t, float64(-55.55), billResponse.Bills[2].Amount)
		assert.Equal(t, person1.ID, billResponse.Bills[3].PersonId)
		assert.Equal(t, float64(55.55), billResponse.Bills[3].Amount)

		// person1
		billResponse, err = GetPendingBills(person1.ID, true, true, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 2)
		assert.Equal(t, billResponse.Count, 2)
		assert.Equal(t, billResponse.FilterPersonId, person1.ID)
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person1.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(-55.55), billResponse.Bills[0].Amount)
		assert.Equal(t, person1.ID, billResponse.Bills[1].PersonId)
		assert.Equal(t, float64(55.55), billResponse.Bills[1].Amount)

		// person2
		billResponse, err = GetPendingBills(person2.ID, true, true, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 2)
		assert.Equal(t, billResponse.Count, 2)
		assert.Equal(t, billResponse.FilterPersonId, person2.ID)
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person2.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(-77.77), billResponse.Bills[0].Amount)
		assert.Equal(t, person2.ID, billResponse.Bills[1].PersonId)
		assert.Equal(t, float64(77.77), billResponse.Bills[1].Amount)

		// to_charge only
		billResponse, err = GetPendingBills(uuid.UUID{}, false, true, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 2)
		assert.Equal(t, billResponse.Count, 2)
		assert.Equal(t, billResponse.FilterPersonId, uuid.UUID{})
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, billResponse.Bills[0].PersonId, person2.ID)
		assert.Equal(t, float64(77.77), billResponse.Bills[0].Amount)
		assert.Equal(t, person1.ID, billResponse.Bills[1].PersonId)
		assert.Equal(t, float64(55.55), billResponse.Bills[1].Amount)

		// to_pay only
		billResponse, err = GetPendingBills(uuid.UUID{}, true, false, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 2)
		assert.Equal(t, billResponse.Count, 2)
		assert.Equal(t, billResponse.FilterPersonId, uuid.UUID{})
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person2.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(-77.77), billResponse.Bills[0].Amount)
		assert.Equal(t, person1.ID, billResponse.Bills[1].PersonId)
		assert.Equal(t, float64(-55.55), billResponse.Bills[1].Amount)

		// person1 to_charge
		billResponse, err = GetPendingBills(person1.ID, false, true, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 1)
		assert.Equal(t, billResponse.Count, 1)
		assert.Equal(t, billResponse.FilterPersonId, person1.ID)
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person1.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(55.55), billResponse.Bills[0].Amount)

		// person1 to_pay
		billResponse, err = GetPendingBills(person1.ID, true, false, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 1)
		assert.Equal(t, billResponse.Count, 1)
		assert.Equal(t, billResponse.FilterPersonId, person1.ID)
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person1.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(-55.55), billResponse.Bills[0].Amount)

		// person2 to_charge
		billResponse, err = GetPendingBills(person2.ID, false, true, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 1)
		assert.Equal(t, billResponse.Count, 1)
		assert.Equal(t, billResponse.FilterPersonId, person2.ID)
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person2.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(77.77), billResponse.Bills[0].Amount)

		// person2 to_pay
		billResponse, err = GetPendingBills(person2.ID, true, false, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 1)
		assert.Equal(t, billResponse.Count, 1)
		assert.Equal(t, billResponse.FilterPersonId, person2.ID)
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person2.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(-77.77), billResponse.Bills[0].Amount)
	})

	emptyBills()

	t.Run("Error when requesting not to pay and not to charge", func(t *testing.T) {
		_, err := GetPendingBills(uuid.UUID{}, false, false, config.Limit, config.Offset)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.BL001, err.Error())
	})

	t.Run("Error when creating transaction with balance = 0", func(t *testing.T) {
		billFields := GenerateBillFields(person1.ID)
		billFields.Amount = 0
		_, err := CreatePendingBill(billFields, uuid.UUID{}, uuid.UUID{})
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.BL002, err.Error())
	})

	t.Run("Error when creating bill with unregistered currency", func(t *testing.T) {
		billFields := GenerateBillFields(person1.ID)
		billFields.Currency = "EEE"
		_, err = CreatePendingBill(billFields, uuid.UUID{}, uuid.UUID{})
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.CU005, err.Error())
	})

	t.Run("Create 14 bills, 10 first random, 2 for person1, 2 for person2, negative and positive balance and get them filtered, paginated", func(t *testing.T) {
		// person2 only has negative balances
		firstBill := Bill{}
		for i := 1; i <= 6; i++ {
			person_id := uuid.UUID{}
			amount := float64(0)
			if i%2 == 0 {
				person_id = person2.ID
				amount = utility.GetRandomNegativeBalance()
			} else {
				person_id = person1.ID
				amount = utility.GetRandomPositiveBalance()
			}
			fields := GenerateBillFields(person_id)
			fields.Amount = amount
			createdBill, err := CreatePendingBill(fields, uuid.UUID{}, uuid.UUID{})
			assert.Nil(t, err)
			if i == 1 {
				firstBill = createdBill
			}
		}

		// person2 only has positive balances
		for i := 1; i <= 6; i++ {
			person_id := uuid.UUID{}
			amount := float64(0)
			if i%2 == 0 {
				person_id = person1.ID
				amount = utility.GetRandomNegativeBalance()
			} else {
				person_id = person2.ID
				amount = utility.GetRandomPositiveBalance()
			}
			fields := GenerateBillFields(person_id)
			fields.Amount = amount
			_, err := CreatePendingBill(fields, uuid.UUID{}, uuid.UUID{})
			assert.Nil(t, err)
		}

		// the last four inserts should be the first ones
		// person1
		billFields := GenerateBillFields(person1.ID)
		billFields.Amount = 55.55
		_, err := CreatePendingBill(billFields, uuid.UUID{}, uuid.UUID{})
		assert.Nil(t, err)

		billFields = GenerateBillFields(person1.ID)
		billFields.Amount = -55.55
		_, err = CreatePendingBill(billFields, uuid.UUID{}, uuid.UUID{})
		assert.Nil(t, err)

		// person2
		billFields = GenerateBillFields(person2.ID)
		billFields.Amount = 77.77
		_, err = CreatePendingBill(billFields, uuid.UUID{}, uuid.UUID{})
		assert.Nil(t, err)

		billFields = GenerateBillFields(person2.ID)
		billFields.Amount = -77.77
		_, err = CreatePendingBill(billFields, uuid.UUID{}, uuid.UUID{})
		assert.Nil(t, err)

		// all of them
		billResponse, err := GetPendingBills(uuid.UUID{}, true, true, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 10)
		assert.Equal(t, 16, billResponse.Count)
		assert.Equal(t, billResponse.FilterPersonId, uuid.UUID{})
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person2.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(-77.77), billResponse.Bills[0].Amount)
		assert.Equal(t, person2.ID, billResponse.Bills[1].PersonId)
		assert.Equal(t, float64(77.77), billResponse.Bills[1].Amount)
		assert.Equal(t, person1.ID, billResponse.Bills[2].PersonId)
		assert.Equal(t, float64(-55.55), billResponse.Bills[2].Amount)
		assert.Equal(t, person1.ID, billResponse.Bills[3].PersonId)
		assert.Equal(t, float64(55.55), billResponse.Bills[3].Amount)

		// second page
		billResponse, err = GetPendingBills(uuid.UUID{}, true, true, config.Limit, 10)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 6)
		assert.Equal(t, 16, billResponse.Count)
		assert.Equal(t, billResponse.FilterPersonId, uuid.UUID{})
		assert.Equal(t, billResponse.Offset, 10)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, firstBill.ID, billResponse.Bills[5].ID)
		assert.Equal(t, firstBill.PersonId, billResponse.Bills[5].PersonId)
		assert.Equal(t, firstBill.Amount, billResponse.Bills[5].Amount)
		assert.Equal(t, firstBill.Currency, billResponse.Bills[5].Currency)
		assert.Equal(t, firstBill.Description, billResponse.Bills[5].Description)

		// person1
		billResponse, err = GetPendingBills(person1.ID, true, true, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 8)
		assert.Equal(t, billResponse.Count, 8)
		assert.Equal(t, billResponse.FilterPersonId, person1.ID)
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person1.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(-55.55), billResponse.Bills[0].Amount)
		assert.Equal(t, person1.ID, billResponse.Bills[1].PersonId)
		assert.Equal(t, float64(55.55), billResponse.Bills[1].Amount)

		// person2
		billResponse, err = GetPendingBills(person2.ID, true, true, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 8)
		assert.Equal(t, billResponse.Count, 8)
		assert.Equal(t, billResponse.FilterPersonId, person2.ID)
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person2.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(-77.77), billResponse.Bills[0].Amount)
		assert.Equal(t, person2.ID, billResponse.Bills[1].PersonId)
		assert.Equal(t, float64(77.77), billResponse.Bills[1].Amount)

		// to_charge only
		billResponse, err = GetPendingBills(uuid.UUID{}, false, true, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 8)
		assert.Equal(t, billResponse.Count, 8)
		assert.Equal(t, billResponse.FilterPersonId, uuid.UUID{})
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, billResponse.Bills[0].PersonId, person2.ID)
		assert.Equal(t, float64(77.77), billResponse.Bills[0].Amount)
		assert.Equal(t, person1.ID, billResponse.Bills[1].PersonId)
		assert.Equal(t, float64(55.55), billResponse.Bills[1].Amount)

		// to_pay only
		billResponse, err = GetPendingBills(uuid.UUID{}, true, false, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 8)
		assert.Equal(t, billResponse.Count, 8)
		assert.Equal(t, billResponse.FilterPersonId, uuid.UUID{})
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person2.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(-77.77), billResponse.Bills[0].Amount)
		assert.Equal(t, person1.ID, billResponse.Bills[1].PersonId)
		assert.Equal(t, float64(-55.55), billResponse.Bills[1].Amount)

		// person1 to_charge
		billResponse, err = GetPendingBills(person1.ID, false, true, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 4)
		assert.Equal(t, billResponse.Count, 4)
		assert.Equal(t, billResponse.FilterPersonId, person1.ID)
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person1.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(55.55), billResponse.Bills[0].Amount)

		// person1 to_pay
		billResponse, err = GetPendingBills(person1.ID, true, false, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 4)
		assert.Equal(t, billResponse.Count, 4)
		assert.Equal(t, billResponse.FilterPersonId, person1.ID)
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person1.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(-55.55), billResponse.Bills[0].Amount)

		// person2 to_charge
		billResponse, err = GetPendingBills(person2.ID, false, true, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 4)
		assert.Equal(t, billResponse.Count, 4)
		assert.Equal(t, billResponse.FilterPersonId, person2.ID)
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person2.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(77.77), billResponse.Bills[0].Amount)

		// person2 to_pay
		billResponse, err = GetPendingBills(person2.ID, true, false, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, billResponse.Bills, 4)
		assert.Equal(t, billResponse.Count, 4)
		assert.Equal(t, billResponse.FilterPersonId, person2.ID)
		assert.Equal(t, billResponse.Offset, config.Offset)
		assert.Equal(t, billResponse.Limit, config.Limit)
		assert.Equal(t, person2.ID, billResponse.Bills[0].PersonId)
		assert.Equal(t, float64(-77.77), billResponse.Bills[0].Amount)
	})

	emptyBills()

	t.Run("Create one bill and get it with single response", func(t *testing.T) {
		billFields := GenerateBillFields(person1.ID)
		newBill, err := CreatePendingBill(billFields, uuid.UUID{}, uuid.UUID{})
		assert.Nil(t, err)
		bill, err := GetOneBill(newBill.ID)
		assert.Nil(t, err)
		assert.Equal(t, newBill.ID, bill.ID)
		assert.Equal(t, newBill.PersonId, bill.PersonId)
		assert.Equal(t, newBill.PersonName, bill.PersonName)
		assert.Equal(t, newBill.Date, bill.Date)
		assert.Equal(t, newBill.Description, bill.Description)
		assert.Equal(t, newBill.Currency, bill.Currency)
		assert.Equal(t, newBill.Amount, bill.Amount)
		assert.Equal(t, newBill.CreatedAt, bill.CreatedAt)
		assert.Equal(t, newBill.UpdatedAt, bill.UpdatedAt)
	})

	emptyBills()

	t.Run("Create closed bill artifitially and get it with single response", func(t *testing.T) {
		billFields := GenerateBillFields(person1.ID)
		newBill, err := createClosedBill(billFields)
		assert.Nil(t, err)
		bill, err := GetOneBill(newBill.ID)
		assert.Nil(t, err)
		assert.Equal(t, newBill.ID, bill.ID)
		assert.Equal(t, newBill.PersonId, bill.PersonId)
		assert.Equal(t, newBill.PersonName, bill.PersonName)
		assert.Equal(t, newBill.Date, bill.Date)
		assert.Equal(t, newBill.Description, bill.Description)
		assert.Equal(t, newBill.Currency, bill.Currency)
		assert.Equal(t, newBill.Amount, bill.Amount)
		assert.Equal(t, newBill.CreatedAt, bill.CreatedAt)
		assert.Equal(t, newBill.UpdatedAt, bill.UpdatedAt)
	})

	t.Run("Error when requesting unexisting bill", func(t *testing.T) {
		randomUUID, err := uuid.NewRandom()
		assert.Nil(t, err)
		_, err = GetOneBill(randomUUID)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB008, err.Error())
	})

	// t.Run("Create one bill and update it", func(t *testing.T) {
	// 	bill, err := CreatePendingBill(GenerateBillFields(person1.ID))
	// 	assert.Nil(t, err)
	// 	updateFields := GenerateBillFields(person1.ID)
	// 	updatedBill, err := UpdatePendingBill(bill.ID, updateFields)
	// 	assert.Nil(t, err)

	// 	assert.Nil(t, err)
	// 	assert.Equal(t, updatedBill.PersonId, updateFields.PersonId)
	// 	assert.Equal(t, updatedBill.Date.Format("2006-01-02"), updateFields.Date.Format("2006-01-02"))
	// 	assert.Equal(t, updatedBill.Description, updateFields.Description)
	// 	// assert.Equal(t, updatedBill.Currency, updateFields.Currency) -> can't update currency
	// 	assert.Equal(t, updatedBill.Amount, updateFields.Amount)

	// 	bill2, err := GetOneBill(bill.ID)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, updatedBill.ID, bill2.ID)
	// 	assert.Equal(t, updatedBill.PersonId, bill2.PersonId)
	// 	assert.Equal(t, updatedBill.PersonName, bill2.PersonName)
	// 	assert.Equal(t, updatedBill.Date, bill2.Date)
	// 	assert.Equal(t, updatedBill.Description, bill2.Description)
	// 	assert.Equal(t, updatedBill.Currency, bill2.Currency)
	// 	assert.Equal(t, updatedBill.Amount, bill2.Amount)
	// 	assert.Equal(t, updatedBill.Pending, bill2.Pending)
	// 	assert.Equal(t, updatedBill.CreatedAt, bill2.CreatedAt)
	// 	assert.Equal(t, updatedBill.UpdatedAt, bill2.UpdatedAt)
	// })

	emptyBills()

	// t.Run("Error when updating unexisting bill", func(t *testing.T) {
	// 	randomUUID, err := uuid.NewRandom()
	// 	assert.Nil(t, err)
	// 	_, err = UpdatePendingBill(randomUUID, GenerateBillFields(person1.ID))
	// 	assert.NotNil(t, err)
	// 	assert.Equal(t, errors_handler.DB001, err.Error())
	// })

}
