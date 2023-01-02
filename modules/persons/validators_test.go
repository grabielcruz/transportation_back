package persons

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckPersonFields(t *testing.T) {
	fields := PersonFields{}
	err := checkPersonFields(fields)
	assert.Equal(t, "Name is required", err.Error())
}
