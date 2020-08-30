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

func TestGenerateSomeRandomString(t *testing.T) {
	size := 10
	n := 5
	randomStrings := make([]string, n)
	for i := 0; i < n; i++ {
		s := GenerateRandomString(size)
		assert.NotContains(t, randomStrings, s)
		randomStrings[i] = s
	}
}
