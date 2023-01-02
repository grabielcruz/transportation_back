package currencies

import (
	"testing"

	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/stretchr/testify/assert"
)

func TestCheckValidCurrency(t *testing.T) {
	goodCurrency := "ABC"
	err := checkValidCurrency(goodCurrency)
	assert.Nil(t, err)

	longCurrency := "ABCD"
	err = checkValidCurrency(longCurrency)
	assert.Equal(t, errors_handler.CU002, err.Error())

	lowerCurrency := "wer"
	err = checkValidCurrency(lowerCurrency)
	assert.Equal(t, errors_handler.CU002, err.Error())

	shortCurrency := "WE"
	err = checkValidCurrency(shortCurrency)
	assert.Equal(t, errors_handler.CU002, err.Error())

	badCurrency := "!@#"
	err = checkValidCurrency(badCurrency)
	assert.Equal(t, errors_handler.CU002, err.Error())

	badCurrency2 := "ABCe"
	err = checkValidCurrency(badCurrency2)
	assert.Equal(t, errors_handler.CU002, err.Error())

	badCurrency3 := "AB C"
	err = checkValidCurrency(badCurrency3)
	assert.Equal(t, errors_handler.CU002, err.Error())
}
