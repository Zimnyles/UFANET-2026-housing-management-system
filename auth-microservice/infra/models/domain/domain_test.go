package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableNames(t *testing.T) {
	assert.Equal(t, "users", (User{}).TableName())
}
