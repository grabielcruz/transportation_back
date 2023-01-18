package transactions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckTransactionFields(t *testing.T) {
	fields := TransactionFields{}
	err := checkTransactionFields(fields)
	assert.Equal(t, "Transaction should have a description", err.Error())
	fields.Description = "asdfasdf asdfas"
	err = checkTransactionFields(fields)
	assert.Equal(t, "Amount should be greater than zero", err.Error())
}
