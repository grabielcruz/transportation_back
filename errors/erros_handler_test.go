package errors_handler

import (
	"os"
	"strings"
	"testing"

	"github.com/grabielcruz/transportation_back/utility"
	"github.com/stretchr/testify/assert"
)

const TestPath = "./errors_log_test.txt"

func TestHandleError(t *testing.T) {
	ResetFile(TestPath)

	count := 100
	// generate an array of 100 errors
	errors := []string{}
	for i := 0; i < count; i++ {
		error_string := utility.GetRandomString(55)
		errors = append(errors, error_string)
	}
	assert.Len(t, errors, count)

	// write errors to errors_log file located in root
	for _, v := range errors {
		WriteErrorToFile(TestPath, v)
	}

	// Read file and check errors are correct
	file, err := os.ReadFile(TestPath)
	assert.Nil(t, err)
	read_errors := strings.Split(string(file), "\n")
	assert.Len(t, read_errors, count+1)
	for i := 0; i < count; i++ {
		assert.Equal(t, read_errors[i], errors[i])
	}

}
