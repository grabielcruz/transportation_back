package bills

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCheckBillFields(t *testing.T) {
	fields := BillFields{}
	err := checkBillFields(fields)
	assert.Equal(t, "Person id should be not zero uuid", err.Error())
	randId, err := uuid.NewRandom()
	assert.Nil(t, err)
	fields.PersonId = randId
	err = checkBillFields(fields)
	assert.Equal(t, "Description is required", err.Error())
	fields.Description = "abc"
	err = checkBillFields(fields)
	assert.Equal(t, "Amount should be greater than zero", err.Error())
	fields.Amount = float64(55)
	err = checkBillFields(fields)
	assert.Equal(t, "Currency code should be 3 upper case letters", err.Error())
	fields.Currency = "ABC"
	err = checkBillFields(fields)
	assert.Nil(t, err)
}
