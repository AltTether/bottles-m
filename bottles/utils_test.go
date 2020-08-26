package bottles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestGenerateRandomString(t *testing.T) {
	size := 10
	s := GenerateRandomString(size)
	assert.Equal(t, size, len(s))
}
